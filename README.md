# CREATE SCHEMA OWNER
```shell
curl --location 'http://localhost:8086/subjects/campaign.campaign_owner_value/versions' \
--header 'Content-Type: application/json' \
--data '{
      "schema": "{\"type\":\"record\",\"name\":\"owner\",\"namespace\":\"campaign.campaign_owner_value\",\"fields\":[{\"name\":\"id\",\"type\":{\"type\":\"string\",\"logicalType\":\"UUID\"}},{\"name\":\"email\",\"type\":\"string\"},{\"name\":\"status\",\"type\":\"string\"},{\"name\":\"user\",\"type\":\"string\"},{\"name\":\"eventTime\",\"type\":{\"type\":\"long\",\"logicalType\":\"timestamp-millis\"}}]}"
}'
```

# CREATE SCHEMA SLUG
```shell
curl --location 'http://localhost:8086/subjects/campaign.campaign_slug_value/versions' \
--header 'Content-Type: application/json' \
--data '{
      "schema": "{\"type\":\"record\",\"name\":\"slug\",\"namespace\":\"campaign.campaign_slug_value\",\"fields\":[{\"name\":\"id\",\"type\":{\"type\":\"string\",\"logicalType\":\"UUID\"}},{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"status\",\"type\":\"string\"},{\"name\":\"cost\",\"type\":\"double\"},{\"name\":\"user\",\"type\":\"string\"},{\"name\":\"eventTime\",\"type\":{\"type\":\"long\",\"logicalType\":\"timestamp-millis\"}}]}"
}'
```

# CREATE SCHEMA REGION
```shell
curl --location 'http://localhost:8086/subjects/campaign.campaign_region_value/versions' \
--header 'Content-Type: application/json' \
--data '{
      "schema": "{\"type\":\"record\",\"name\":\"region\",\"namespace\":\"campaign.campaign_region_value\",\"fields\":[{\"name\":\"id\",\"type\":{\"type\":\"string\",\"logicalType\":\"UUID\"}},{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"status\",\"type\":\"string\"},{\"name\":\"lat\",\"type\":\"double\"},{\"name\":\"long\",\"type\":\"double\"},{\"name\":\"cost\",\"type\":\"double\"},{\"name\":\"user\",\"type\":\"string\"},{\"name\":\"eventTime\",\"type\":{\"type\":\"long\",\"logicalType\":\"timestamp-millis\"}}]}"
}'
```

# CREATE SCHEMA MERCHANT
```shell

```

# CREATE SCHEMA CAMPAIGN
```shell

```