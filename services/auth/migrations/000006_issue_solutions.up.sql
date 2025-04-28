CREATE TABLE issue_solutions (
                           solution_id BIGSERIAL PRIMARY KEY,
                           assignment_id BIGINT NOT NULL REFERENCES issue_assignments(assignment_id) ON DELETE CASCADE,
                           solution_text TEXT NOT NULL,
                           status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'checked'))
);

ALTER TABLE issue_solutions ADD CONSTRAINT unique_assignment UNIQUE (assignment_id);
