package repository

import (
	"company/internal/domain/entities"
	"company/internal/storage"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
)

func (s *Storage) CreateIssue(ctx context.Context, installationID int64, title string, body string) (int64, error) {
	const op = "repository.issue.AddIssue"

	var id int64

	query := `INSERT INTO issues (installation_id, title, body) VALUES ($1, $2, $3) RETURNING id`

	args := []any{installationID, title, body}
	err := s.Db.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	return id, nil
}

func (s *Storage) GetIssue(ctx context.Context, issueID int64, developerID *int64) (*entities.Issue, error) {
	const op = "repository.issue.GetIssue"
	var queryStr string

	if developerID == nil {
		queryStr = `SELECT id, title, body FROM issues WHERE id = $1`
	} else {
		queryStr = `
			SELECT 
				i.id,
				i.title,
				i.body,
				a.status AS assignment_status,
				s.solution_text,
				s.status AS solution_status
			FROM issues i
			LEFT JOIN issue_assignments a ON i.id = a.issue_id AND a.developer_id = $2
			LEFT JOIN issue_solutions s ON a.assignment_id = s.assignment_id
			WHERE i.id = $1
		`
	}

	stmt, err := s.Db.Prepare(queryStr)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	var issue entities.Issue
	if developerID == nil {
		err = stmt.QueryRowContext(ctx, issueID).Scan(&issue.ID, &issue.Title, &issue.Body)
	} else {
		err = stmt.QueryRowContext(ctx, issueID, *developerID).Scan(
			&issue.ID,
			&issue.Title,
			&issue.Body,
			&issue.AssignmentStatus,
			&issue.SolutionText,
			&issue.SolutionStatus,
		)
	}

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, storage.ErrNoRecordFound
		default:
			return nil, fmt.Errorf("%s:%w", op, err)
		}
	}

	return &issue, nil
}

func (s *Storage) GetAllCompanyIssues(ctx context.Context, id int64) ([]*entities.Issue, error) {
	const op = "repository.issue.GetAllIssues"

	query := "SELECT i.id, i.title, i.body FROM issues i INNER JOIN companies c ON i.installation_id = c.installation_id WHERE c.company_id = $1"

	rows, err := s.Db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	//totalRecords := 0
	var issues []*entities.Issue

	for rows.Next() {
		issue := &entities.Issue{}

		err := rows.Scan(
			&issue.ID,
			&issue.Title,
			&issue.Body,
		)

		if err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		issues = append(issues, issue)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	//if totalRecords == 0 {
	//	return nil, storage.ErrNoRecordFound
	//}

	return issues, nil
}

func (s *Storage) CreateAssignment(ctx context.Context, issueID, developerID int64) (int64, error) {
	const op = "repository.issue.CreateAssignment"

	var id int64

	query := `INSERT INTO issue_assignments (issue_id, developer_id) VALUES ($1, $2) RETURNING assignment_id`

	args := []any{issueID, developerID}
	err := s.Db.QueryRowContext(ctx, query, args...).Scan(&id)
	if err != nil {
		var pqErr *pq.Error
		switch {
		case errors.As(err, &pqErr) && pqErr.Code == "23505":
			return 0, storage.AlreadyExists

		default:
			return 0, fmt.Errorf("%s:%w", op, err)
		}
	}

	return id, nil
}

func (s *Storage) CreateSolution(ctx context.Context, assignmentID int64, solution string) (int64, error) {
	const op = "repository.issue.CreateSolution"

	tx, err := s.Db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	// TODO: Wrap
	defer tx.Rollback()

	var solutionID int64
	err = tx.QueryRowContext(ctx, `
		INSERT INTO issue_solutions (assignment_id, solution_text)
		VALUES ($1, $2) RETURNING solution_id
	`, assignmentID, solution).Scan(&solutionID)

	if err != nil {
		var pqErr *pq.Error
		switch {
		case errors.As(err, &pqErr) && pqErr.Code == "23505":
			return 0, storage.AlreadyExists
		default:
			return 0, fmt.Errorf("%s:%w", op, err)
		}
	}

	_, err = tx.ExecContext(ctx, `
		UPDATE issue_assignments 
		SET status = 'completed', completed_at = NOW() 
		WHERE assignment_id = $1
	`, assignmentID)
	if err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	return solutionID, nil

}

func (s *Storage) GetAllIssueSolutions(ctx context.Context, id int64) ([]*entities.Solution, error) {
	const op = "repository.issue.GetAllSolutions"

	query := `
        SELECT 
            s.solution_id, 
            s.assignment_id, 
            s.solution_text, 
            s.status,
            a.assigned_at,
            a.completed_at
        FROM 
            issue_solutions s
        JOIN 
            issue_assignments a ON s.assignment_id = a.assignment_id
        WHERE 
            a.issue_id = $1
    `

	rows, err := s.Db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {

		}
	}(rows)

	//totalRecords := 0
	var solutions []*entities.Solution

	for rows.Next() {
		solution := &entities.Solution{}

		err := rows.Scan(
			&solution.ID,
			&solution.AssignmentID,
			&solution.SolutionText,
			&solution.Status,
			&solution.AssignedAt,
			&solution.CompletedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("%s:%w", op, err)
		}

		solutions = append(solutions, solution)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s:%w", op, err)
	}

	if len(solutions) == 0 {
		return nil, storage.ErrNoRecordFound
	}

	return solutions, nil
}

func (s *Storage) GetSolution(ctx context.Context, issueID int64, solutionID int64) (*entities.Solution, error) {
	const op = "repository.issue.GetSolution"

	query := `
        SELECT 
            s.solution_id, 
            s.assignment_id, 
            s.solution_text, 
            s.status,
            a.assigned_at,
            a.completed_at
        FROM 
            issue_solutions s
        JOIN 
            issue_assignments a ON s.assignment_id = a.assignment_id
        WHERE 
            a.issue_id = $1 AND s.solution_id = $2
    `

	solution := &entities.Solution{}
	err := s.Db.QueryRowContext(ctx, query, issueID, solutionID).Scan(
		&solution.ID,
		&solution.AssignmentID,
		&solution.SolutionText,
		&solution.Status,
		&solution.AssignedAt,
		&solution.CompletedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, storage.ErrNoRecordFound
		default:
			return nil, fmt.Errorf("%s:%w", op, err)
		}
	}

	return solution, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////

func (s *Storage) GetAssignmentID(ctx context.Context, issueID, developerID int64) (int64, error) {
	const op = "repository.issue.GetAssignmentID"

	var assignmentID int64

	err := s.Db.QueryRowContext(ctx, `
		SELECT assignment_id FROM issue_assignments 
		WHERE issue_id = $1 AND developer_id = $2
	`, issueID, developerID).Scan(&assignmentID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, storage.ErrNoRecordFound
		}
		return 0, fmt.Errorf("%s:%w", op, err)
	}
	return assignmentID, nil
}

func (s *Storage) GetCompanyIDByIssueID(ctx context.Context, issueID int64) (int64, error) {
	const op = "repository.issue.GetCompanyIDByIssueID"

	var companyID int64

	err := s.Db.QueryRowContext(ctx, `
		SELECT c.company_id
		FROM issues i
		JOIN companies c ON i.installation_id = c.installation_id
		WHERE i.id = $1
	`, issueID).Scan(&companyID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, storage.ErrNoRecordFound
		}
		return 0, fmt.Errorf("%s:%w", op, err)
	}

	return companyID, nil
}
