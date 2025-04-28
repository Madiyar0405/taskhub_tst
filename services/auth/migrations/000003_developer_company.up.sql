CREATE TABLE companies (
                           company_id BIGSERIAL PRIMARY KEY,
                           company_name VARCHAR(255) NOT NULL,
                           description TEXT,
                           website VARCHAR(255),
                           logo VARCHAR(255),
                           installation_id BIGINT UNIQUE,
                           created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
                           updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW()
--                            project_key BIGSERIAL UNIQUE
);

CREATE TABLE developers (
                            developer_id BIGSERIAL PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
                            github_id BIGINT UNIQUE,
                            bio TEXT,
                            github_username VARCHAR(255),
                            avatar_url VARCHAR(255),
                            cv_url VARCHAR(255)
);