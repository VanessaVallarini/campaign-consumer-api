CREATE TABLE spent(
    id UUID PRIMARY KEY NOT NULL,
    campaign_id  UUID NOT NULL,
    bucket VARCHAR(50) NOT NULL,
    total_spent DECIMAL(14,2) NOT NULL,
    total_clicks integer NOT NULL,
    total_impressions integer NOT NULL,
    FOREIGN KEY (campaign_id) REFERENCES campaign(id)
);