CREATE TYPE status_service_enum AS ENUM ('pending', 'online', 'offline');

CREATE TABLE services (
  id UUID PRIMARY KEY,
  name TEXT NOT NULL,
  url_address TEXT NOT NULL,
  status status_service_enum NOT NULL DEFAULT 'pending',
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ
);
