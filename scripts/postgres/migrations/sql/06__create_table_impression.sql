CREATE TABLE impression(
    campaign_id  UUID PRIMARY KEY NOT NULL,
    lat FLOAT,
    long FLOAT,
    created_at TIMESTAMP WITH TIME ZONE,
    user_id VARCHAR(60),
    FOREIGN KEY (campaign_id) REFERENCES campaign(id)
);