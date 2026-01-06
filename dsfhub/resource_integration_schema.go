package dsfhub

var requiredCiphertrustFieldsJson = `{
	"id",
	"type",
	"hostname",
	"port",
	"cm_name",
	"auth_method",
	"is_load_balancer",
}`

var requiredClassificationFieldsJson = `{
	"id",
	"type",
	"storage_details": {
		"storage_type",
		"s3_bucket_configuration: {
			"bucket_name",
			"cloud_name",
			"access_key_id",
			"secret_access_key",
			"aws_region"
		}
	},
	"database_details: {
		"database_type",
		"mongo_configuration: {
			"db_name",
			"connection_string",
		}
	}
}`
