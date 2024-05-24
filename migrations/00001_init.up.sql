CREATE TABLE IF NOT EXISTS orders (
	id UUID PRIMARY KEY DEFAULT get_random_uuid(),
	orderItem JSONB NOT NULL,
)
