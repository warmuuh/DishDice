CREATE TABLE IF NOT EXISTS registration_tickets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    token VARCHAR(255) UNIQUE NOT NULL,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    used_by UUID REFERENCES users(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    used_at TIMESTAMP,
    is_used BOOLEAN DEFAULT false
);

CREATE INDEX idx_tickets_token ON registration_tickets(token);
CREATE INDEX idx_tickets_created_by ON registration_tickets(created_by);
CREATE INDEX idx_tickets_used_status ON registration_tickets(is_used, expires_at);
