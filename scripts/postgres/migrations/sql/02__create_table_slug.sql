CREATE TABLE IF NOT EXISTS slug(
    id  UUID PRIMARY KEY NOT NULL,
    name   VARCHAR(50) NOT NULL,
    active BOOLEAN NOT NULL,
    cost DECIMAL(5,2) NOT NULL,
    created_by VARCHAR(60),
    updated_by VARCHAR(60),
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);