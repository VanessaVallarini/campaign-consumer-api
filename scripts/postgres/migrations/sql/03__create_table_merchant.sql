CREATE TABLE merchant(
    id  UUID PRIMARY KEY NOT NULL,
    owner_id  UUID NOT NULL,
    slugs UUID[] NOT NULL,
    name   VARCHAR(50) NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    created_by VARCHAR(60),
    updated_by VARCHAR(60),
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (owner_id) REFERENCES owner(id)
);