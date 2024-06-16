# CREATE SCHEMA OWNER
```shell
curl --location 'http://localhost:8086/subjects/campaign.campaign_owner_value/versions' \
--header 'Content-Type: application/json' \
--data '{
        "schema": "{\"type\":\"record\",\"name\":\"owner\",\"namespace\":\"campaign.campaign_owner_value\",\"fields\":[{\"name\":\"id\",\"type\":{\"type\":\"string\",\"logicalType\":\"UUID\"}},{\"name\":\"email\",\"type\":\"string\"},{\"name\":\"status\",\"type\":\"string\"},{\"name\":\"created_by\",\"type\":\"string\"},{\"name\":\"updated_by\",\"type\":\"string\"},{\"name\":\"created_at\",\"type\":{\"type\":\"long\",\"logicalType\":\"timestamp-millis\"}},{\"name\":\"updated_at\",\"type\":{\"type\":\"long\",\"logicalType\":\"timestamp-millis\"}}]}"

}'
```

# CREATE SCHEMA SLUG
```shell
```

# CREATE SCHEMA REGION
```shell
```

# CREATE SCHEMA MERCHANT
```shell
```

# CREATE SCHEMA CAMPAIGN
```shell
```

# CREATE SCHEMA CLICK IMPRESSION
```shell
```