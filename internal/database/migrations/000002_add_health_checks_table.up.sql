CREATE TABLE health_checks (
    id UUID PRIMARY KEY,
    
    service_id UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    
    status VARCHAR(50) NOT NULL,
    
    response_time_ms INTEGER NOT NULL,
    
    http_status_code INTEGER,
    
    error_message TEXT,
    
    checked_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_health_checks_service_id ON health_checks(service_id);