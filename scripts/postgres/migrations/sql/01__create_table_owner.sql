CREATE TABLE owner(
    id  UUID PRIMARY KEY NOT NULL,
    email   VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE,
    active BOOLEAN DEFAULT TRUE
);