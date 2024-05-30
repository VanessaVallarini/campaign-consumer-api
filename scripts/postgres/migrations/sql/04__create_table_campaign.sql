CREATE TABLE campaign(
    id  UUID PRIMARY KEY NOT NULL,
    merchant_id  UUID NOT NULL,
    active BOOLEAN DEFAULT TRUE,
    lat FLOAT NOT NULL,
    long FLOAT NOT NULL,
    created_by VARCHAR(60),
    updated_by VARCHAR(60),
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE,
    FOREIGN KEY (merchant_id) REFERENCES merchant(id)
);