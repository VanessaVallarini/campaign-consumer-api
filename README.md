# CREATE SCHEMA OWNER
```shell
curl --location 'http://localhost:8086/subjects/campaign.campaign_owner_value/versions' \
--header 'Content-Type: application/json' \
--data '{
    "schema": "{\"type\":\"record\",\"name\":\"owner\",\"namespace\":\"campaign.campaign_owner_value\",\"fields\":[{\"name\":\"id\",\"type\":\"string\"},{\"name\":\"email\",\"type\":\"string\"},{\"name\":\"active\",\"type\":\"string\"},{\"name\":\"created_by\",\"type\":\"string\"},{\"name\":\"updated_by\",\"type\":\"string\"}]}"
}'
```

# SEND MESSAGE OWNER 
```shell
{
	"id": "7631fcd2-5722-47d9-86af-761b4ca16644",
	"email": "van@teste.com",
	"active": "true",
	"created_by": "van",
	"updated_by": "van"
}
```

# CREATE SCHEMA SLUG
```shell
curl --location 'http://localhost:8086/subjects/campaign.campaign_slug_value/versions' \
--header 'Content-Type: application/json' \
--data '{
    "schema": "{\"type\":\"record\",\"name\":\"slug\",\"namespace\":\"campaign.campaign_slug_value\",\"fields\":[{\"name\":\"id\",\"type\":\"string\"},{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"active\",\"type\":\"string\"},{\"name\":\"cost\",\"type\":\"double\"},{\"name\":\"created_by\",\"type\":\"string\"},{\"name\":\"updated_by\",\"type\":\"string\"}]}"
}'
```

# SEND MESSAGE OWNER 
```shell
{
	"id": "550146ad-eaff-46e7-8388-b8f948740903",
	"name": "van",
	"active": "true",
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
    "schema": "{\"type\":\"record\",\"name\":\"region\",\"namespace\":\"campaign.campaign_region_value\",\"fields\":[{\"name\":\"id\",\"type\":\"string\"},{\"name\":\"name\",\"type\":\"string\"},{\"name\":\"active\",\"type\":\"string\"},{\"name\":\"lat\",\"type\":\"double\"},{\"name\":\"long\",\"type\":\"double\"},{\"name\":\"cost\",\"type\":\"double\"},{\"name\":\"created_by\",\"type\":\"string\"},{\"name\":\"updated_by\",\"type\":\"string\"}]}"
}'
```

# SEND MESSAGE REGION 
```shell
{
	"id": "550146ad-eaff-46e7-8388-b8f948740903",
	"name": "Londrina",
	"active": "true",
	"lat": ,
	"long": ,
	"cost": 0.55,
	"created_by": "van",
	"updated_by": "van"
}
```