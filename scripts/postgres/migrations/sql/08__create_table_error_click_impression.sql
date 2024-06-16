CREATE TABLE IF NOT EXISTS error_click_impression (
    id  UUID PRIMARY KEY NOT NULL,
    campaign_id  UUID NOT NULL,
    slug_id  UUID NOT NULL,
    user_id UUID NOT NULL,
    event_type VARCHAR(30) NOT NULL,
    lat FLOAT NOT NULL,
    long FLOAT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    FOREIGN KEY (campaign_id) REFERENCES campaign(id),
    FOREIGN KEY (slug_id) REFERENCES slug(id)
);

CREATE INDEX error_click_impression_id ON campaign_consumer_api.error_click_impression USING btree (id);
CREATE INDEX error_click_impression_campaign_id ON campaign_consumer_api.error_click_impression USING btree (campaign_id);
CREATE INDEX error_click_impression_slug_id ON campaign_consumer_api.error_click_impression USING btree (slug_id);