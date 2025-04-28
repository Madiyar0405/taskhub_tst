CREATE TABLE issue_assignments (
                                   assignment_id BIGSERIAL PRIMARY KEY,
                                   issue_id BIGINT NOT NULL REFERENCES issues(id) ON DELETE CASCADE,
                                   developer_id BIGINT NOT NULL REFERENCES developers(developer_id) ON DELETE CASCADE,
                                   status VARCHAR(20) DEFAULT 'in_progress' CHECK (status IN ('in_progress', 'completed')),
                                   assigned_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
                                   completed_at TIMESTAMP(0) WITH TIME ZONE,
                                   UNIQUE(issue_id, developer_id)
);
