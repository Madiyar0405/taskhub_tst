CREATE TABLE otps (
                      id BIGSERIAL PRIMARY KEY,
                      email CITEXT NOT NULL,
                      otp_code TEXT NOT NULL,
                      created_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
                      expires_at TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT (NOW() + INTERVAL '30 minutes'),
                      FOREIGN KEY (email) REFERENCES users (email) ON DELETE CASCADE
);
