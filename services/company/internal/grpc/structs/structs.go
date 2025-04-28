package structs

type CompanyRequest struct {
	CompanyID int64 `validate:"required,min=0"`
}

type CompanyIntegrationRequest struct {
	CompanyID int64 `validate:"required,min=0"`
}

type IssueRequest struct {
	IssueID int64 `validate:"required,min=0"`
}

type AssignDeveloperRequest struct {
	IssueID int64 `validate:"required,min=0"`
}

type SubmitSolutionRequest struct {
	IssueID  int64  `validate:"required,min=0"`
	Solution string `validate:"required,min=8"`
}

type IssueSolutionsRequest struct {
	IssueID int64 `validate:"required,min=0"`
}

type IssueSolutionRequest struct {
	IssueID    int64 `validate:"required,min=0"`
	SolutionID int64 `validate:"required,min=0"`
}
