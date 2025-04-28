package webhook

import (
	"company/internal/grpc/company"
	"context"
	"fmt"
	"github.com/google/go-github/v50/github"
	"log/slog"
	"net/http"
)

type Handler struct {
	log            *slog.Logger
	companyService company.Company
	issueService   company.Issue
	webhookSecret  string
}

func NewHandler(log *slog.Logger, companyService company.Company, issueService company.Issue, webhookSecret string) *Handler {
	return &Handler{
		log:            log,
		companyService: companyService,
		issueService:   issueService,
		webhookSecret:  webhookSecret,
	}
}

func (h *Handler) HandleWebhook(w http.ResponseWriter, r *http.Request) {

	log := h.log.With(slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("event", r.Header.Get("X-GitHub-Event")))

	payload, err := github.ValidatePayload(r, []byte(h.webhookSecret))
	if err != nil {
		log.Error("Invalid webhook signature")
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	event, err := github.ParseWebHook(github.WebHookType(r), payload)
	if err != nil {
		log.Info("Failed to parse webhook")
		http.Error(w, "Failed to parse webhook", http.StatusBadRequest)
		return
	}

	switch e := event.(type) {
	case *github.InstallationEvent:
		log.Info("Processing installation event")
		h.handleInstallationEvent(e, w)
	case *github.IssuesEvent:
		log.Info("Processing issue event")
		h.handleIssueEvent(e, w)
	default:
		log.Info("Ignoring unsupported event type")
		w.WriteHeader(http.StatusOK)
	}

	log.Info("successfully")
}

func (h *Handler) handleIssueEvent(event *github.IssuesEvent, w http.ResponseWriter) {
	log := h.log.With(
		slog.String("action", event.GetAction()),
		slog.Int64("issue_id", event.GetIssue().GetID()),
		slog.Int64("installation_id", event.GetInstallation().GetID()),
	)

	log.Info("Processing issue event")

	if event.GetAction() != "opened" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if event.GetIssue().IsPullRequest() {
		log.Info("Skipping pull request")
		w.WriteHeader(http.StatusOK)
		return
	}

	installationID := event.GetInstallation().GetID()

	issueTitle := event.Issue.GetTitle()
	bodyTitle := event.Issue.GetBody()

	_, err := h.issueService.AddIssue(context.Background(), installationID, issueTitle, bodyTitle)
	if err != nil {
		log.Error("Failed to process issue", slog.String("error", err.Error()))
		http.Error(w, "Failed to process issue", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) handleInstallationEvent(event *github.InstallationEvent, w http.ResponseWriter) {

	log := h.log.With(
		slog.String("action", event.GetAction()),
		slog.Int64("installation_id", event.Installation.GetID()),
	)

	log.Info("Handling installation event")

	if *event.Action != "created" {
		log.Info("Ignoring non-creation event")
		w.WriteHeader(http.StatusOK)
		return
	}

	installationID := *event.Installation.ID
	account := event.Installation.Account

	//if *account.Type != "Organization" {
	//	log.Info("Ignoring non-organization installation")
	//	w.WriteHeader(http.StatusOK)
	//	return
	//}

	companyName := *account.Login
	if account.Name != nil && *account.Name != "" {
		companyName = *account.Name
	}

	logoURL := ""
	if account.AvatarURL != nil {
		logoURL = *account.AvatarURL
	}

	log = log.With(
		slog.String("company_name", companyName),
		slog.String("logo_url", logoURL),
	)

	log.Info("Attempting to save company")

	companyID, err := h.companyService.AddCompany(context.Background(), installationID, companyName, logoURL)
	if err != nil {
		log.Error("Failed to save company", slog.String("error", err.Error()))
		http.Error(w, fmt.Sprintf("Failed to save company: %v", err), http.StatusInternalServerError)
		return
	}

	log.Info("Company saved successfully", slog.Int64("company_id", companyID))
	w.WriteHeader(http.StatusOK)
}
