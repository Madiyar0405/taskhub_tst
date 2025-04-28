CREATE TABLE users(
                          id BIGSERIAL PRIMARY KEY,
                          username text UNIQUE NOT NULL,
                          email citext UNIQUE NOT NULL,
                          role VARCHAR(20) DEFAULT 'developer' CHECK (role IN ('developer', 'company', 'admin')),
                          created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
                          updated_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
                          password_hash bytea NOT NULL,
                          activated bool NOT NULL
                  );