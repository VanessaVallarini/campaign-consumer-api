# CREATE SCHEMA OWNER
```shell
curl --location 'http://localhost:8086/subjects/campaign.campaign_owner_value/versions' \
--header 'Content-Type: application/json' \
--data '{
    "schema": "{\"type\":\"record\",\"name\":\"owner\",\"namespace\":\"campaign.campaign_owner_value\",\"fields\":[{\"name\":\"id\",\"type\":\"string\"},{\"name\":\"email\",\"type\":\"string\"},{\"name\":\"status\",\"type\":\"string\"},{\"name\":\"created_by\",\"type\":\"string\"},{\"name\":\"updated_by\",\"type\":\"string\"}]}"
}'
```

# SEND MESSAGE OWNER 
```shell
{
	"id": "7631fcd2-5722-47d9-86af-761b4ca16644",
	"email": "van@teste.com",
	"status": "ACTIVE",
	"created_by": "van",
	"updated_by": "van"
}
```

# CREATE SCHEMA SLUG
```shell
curl --location 'http://localhost:8086/subjects/campaign.campaign_slug_value/versions' \
--header 'Content-Type: application/json' \
--data '{
    "schema": "{\"type\":\"record\",\"name\":\"slug\",\"namespace\":\"campaign.campaign_slug_value\",\"fields\":[{\"name\":\"id\",\"type\":\"string\"},{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"status\",\"type\":\"string\"},{\"name\":\"cost\",\"type\":\"double\"},{\"name\":\"created_by\",\"type\":\"string\"},{\"name\":\"updated_by\",\"type\":\"string\"}]}"
}'
```

# SEND MESSAGE SLUG 
```shell
{
	"id": "550146ad-eaff-46e7-8388-b8f948740903",
	"name": "pizza",
	"status": "ACTIVE",
	"cost": 0.55,
	"created_by": "van",
	"updated_by": "van"
}
```

# CREATE SCHEMA REGION
```shell
curl --location 'http://localhost:8086/subjects/campaign.campaign_region_value/versions' \
--header 'Content-Type: application/json' \
--data '{
    "schema": "{\"type\":\"record\",\"name\":\"region\",\"namespace\":\"campaign.campaign_region_value\",\"fields\":[{\"name\":\"id\",\"type\":\"string\"},{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"status\",\"type\":\"string\"},{\"name\":\"lat\",\"type\":\"double\"},{\"name\":\"long\",\"type\":\"double\"},{\"name\":\"cost\",\"type\":\"double\"},{\"name\":\"created_by\",\"type\":\"string\"},{\"name\":\"updated_by\",\"type\":\"string\"}]}"
}'
```

# SEND MESSAGE REGION 
```shell
{
	"id": "e981913f-0c2a-4149-8f0b-548c14433d23",
	"name": "Londrina",
	"status": "ACTIVE",
	"lat": -23.3212795,
	"long": -51.165763,
	"cost": 0.55,
	"created_by": "van",
	"updated_by": "van"
}
```

# CREATE SCHEMA MERCHANT
```shell
curl --location 'http://localhost:8086/subjects/campaign.campaign_merchant_value/versions' \
--header 'Content-Type: application/json' \
--data '{
    "schema": "{\"type\":\"record\",\"name\":\"merchant\",\"namespace\":\"campaign.campaign_merchant_value\",\"fields\":[{\"name\":\"id\",\"type\":\"string\"},{\"name\":\"owner_id\",\"type\":\"string\"},{\"name\":\"region_id\",\"type\":\"string\"},{\"name\":\"slugs\",\"type\":{\"type\":\"array\",\"items\":\"string\"}},{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"status\",\"type\":\"string\"},{\"name\":\"created_by\",\"type\":\"string\"},{\"name\":\"updated_by\",\"type\":\"string\"}]}]}"
}'
```

# SEND MESSAGE MERCHANT 
```shell
{
	"id": "ed48d4be-2e8e-4f58-82af-b568ec3e0097",
	"owner_id": "7631fcd2-5722-47d9-86af-761b4ca16644",
	"region_id": "e981913f-0c2a-4149-8f0b-548c14433d23",
	"slugs": ["550146ad-eaff-46e7-8388-b8f948740903", "adb1363e-3f91-4683-9844-5f937b2b31d5"],
	"name": "Pastelaria Promo",
	"status": "ACTIVE",
	"created_by": "van",
	"updated_by": "van"
}
```