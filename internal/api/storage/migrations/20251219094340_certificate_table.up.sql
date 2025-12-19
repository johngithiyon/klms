CREATE TABLE certificate_info (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    name TEXT NOT NULL,
    issued_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_user
        FOREIGN KEY (username)
        REFERENCES users(username)
        ON DELETE CASCADE,
    CONSTRAINT unique_certificate
        UNIQUE (username)
);
