CREATE TABLE IF NOT EXISTS ledger (
    id  UUID PRIMARY KEY NOT NULL,
    campaign_id  UUID NOT NULL,
    spent_id  UUID NOT NULL,
    slug_id UUID NOT NULL,
    user_id UUID NOT NULL,
    event_type VARCHAR(30) NOT NULL,
    cost DECIMAL(5,2) NOT NULL,
    lat FLOAT NOT NULL,
    long FLOAT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (campaign_id) REFERENCES campaign(id),
    FOREIGN KEY (spent_id) REFERENCES spent(id),
    FOREIGN KEY (slug_id) REFERENCES slug(id)
);
--slugId