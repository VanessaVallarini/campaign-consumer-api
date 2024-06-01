CREATE TABLE IF NOT EXISTS region(
    id  UUID PRIMARY KEY NOT NULL,
    name   VARCHAR(50) NOT NULL,
    active BOOLEAN NOT NULL,
    lat DECIMAL(10,6) NOT NULL,
    long DECIMAL(10,6) NOT null,
    cost DECIMAL(5,2) NOT NULL,
    created_by VARCHAR(60),
    updated_by VARCHAR(60),
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE
);