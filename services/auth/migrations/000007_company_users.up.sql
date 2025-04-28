CREATE TABLE company_users (
                               id BIGSERIAL PRIMARY KEY,
                               user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                               company_id BIGINT NOT NULL REFERENCES companies(company_id) ON DELETE CASCADE,
                               UNIQUE(user_id)
);
