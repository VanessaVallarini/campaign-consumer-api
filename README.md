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
curl --location 'http://localhost:8086/subjects/campaign.campaign_slug_value/versions' \
--header 'Content-Type: application/json' \
--data '{
        "schema": "{\"type\":\"record\",\"name\":\"slug\",\"namespace\":\"campaign.campaign_slug_value\",\"fields\":[{\"name\":\"id\",\"type\":{\"type\":\"string\",\"logicalType\":\"UUID\"}},{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"status\",\"type\":\"string\"},{\"name\":\"cost\",\"type\":\"double\"},{\"name\":\"created_by\",\"type\":\"string\"},{\"name\":\"updated_by\",\"type\":\"string\"},{\"name\":\"created_at\",\"type\":{\"type\":\"long\",\"logicalType\":\"timestamp-millis\"}},{\"name\":\"updated_at\",\"type\":{\"type\":\"long\",\"logicalType\":\"timestamp-millis\"}}]}"

}'
```

# CREATE SCHEMA REGION
```shell
curl --location 'http://localhost:8086/subjects/campaign.campaign_region_value/versions' \
--header 'Content-Type: application/json' \
--data '{
        "schema": "{\"type\":\"record\",\"name\":\"slug\",\"namespace\":\"campaign.campaign_slug_value\",\"fields\":[{\"name\":\"id\",\"type\":{\"type\":\"string\",\"logicalType\":\"UUID\"}},{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"status\",\"type\":\"string\"},{\"name\":\"lat\",\"type\":\"double\"},{\"name\":\"long\",\"type\":\"double\"},{\"name\":\"cost\",\"type\":\"double\"},{\"name\":\"created_by\",\"type\":\"string\"},{\"name\":\"updated_by\",\"type\":\"string\"},{\"name\":\"created_at\",\"type\":{\"type\":\"long\",\"logicalType\":\"timestamp-millis\"}},{\"name\":\"updated_at\",\"type\":{\"type\":\"long\",\"logicalType\":\"timestamp-millis\"}}]}"
}'
```

# CREATE SCHEMA MERCHANT
```shell
curl --location 'http://localhost:8086/subjects/campaign.campaign_merchant_value/versions' \
--header 'Content-Type: application/json' \
--data '{
      "schema": "{\"type\":\"record\",\"name\":\"merchant\",\"namespace\":\"campaign.campaign_merchant_value\",\"fields\":[{\"name\":\"id\",\"type\":{\"type\":\"string\",\"logicalType\":\"UUID\"}},{\"name\":\"owner_id\",\"type\":{\"type\":\"string\",\"logicalType\":\"UUID\"}},{\"name\":\"region_id\",\"type\":{\"type\":\"string\",\"logicalType\":\"UUID\"}},{\"name\":\"slugs\",\"type\":{\"type\":\"array\",\"type\":\"string\",\"logicalType\":\"UUID\"}},{\"name\":\"status\",\"type\":\"string\"},{\"name\":\"created_by\",\"type\":\"string\"},{\"name\":\"updated_by\",\"type\":\"string\"},{\"name\":\"created_at\",\"type\":{\"type\":\"long\",\"logicalType\":\"timestamp-millis\"}},{\"name\":\"updated_at\",\"type\":{\"type\":\"long\",\"logicalType\":\"timestamp-millis\"}}]}"
}'
```

# CREATE SCHEMA CAMPAIGN
```shell
```

# CREATE SCHEMA CLICK IMPRESSION
```shell
```