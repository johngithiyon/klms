CREATE TABLE pending_notifications (
    username   VARCHAR(255) PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_pending_notifications_username
ON pending_notifications (username);
