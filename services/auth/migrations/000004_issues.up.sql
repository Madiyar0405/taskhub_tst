CREATE TABLE issues (
                        id BIGSERIAL PRIMARY KEY,
                        installation_id BIGINT REFERENCES companies(installation_id),
                        title TEXT NOT NULL,
                        body TEXT
);