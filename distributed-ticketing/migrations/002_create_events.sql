CREATE TABLE events (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    total_tickets INT NOT NULL,
    available_tickets INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);