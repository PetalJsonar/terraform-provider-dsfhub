package dsfhub

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceClassification() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClassificationCreateContext,
		ReadContext:   resourceClassificationReadContext,
		UpdateContext: resourceClassificationUpdateContext,
		DeleteContext: resourceClassificationDeleteContext,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "ID of the classification.",
				Optional:    true,
				Computed:   true,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the classification.",
				Optional:    true,
				Computed: true,
				Default: "Data discovery and classification service, used to classify all data within your organization.",
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Type of the classification.",
				Computed:    true,
				Default: "Classification service",
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Status of the classification.",
				Optional:    true,
				Default: "N/A",
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "User-friendly name of the classification, defined by user.",
				Required:    true,
			},
			"last_status_update": {
				Type:        schema.TypeString,
				Description: "Timestamp of the last status update.",
				Optional:    true,
				Default: nil,
			},
			"database_details": {
				Type:        schema.TypeSet,
				Description: "Database details of the classification.",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database_type": {
							Type:        schema.TypeString,
							Description: "Name of the database.",
							Computed:    true,
							Default:    "MongoDB",
						},
						"mongo_configuration": {
							Type:        schema.TypeString,
							Description: "Connection string for MongoDB.",
							Optional:    true,
							Required:    false,
						},
					},
				},
			},
			"storage_details": {
				Type:        schema.TypeSet,
				Description: "Storage details of the classification.",
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_type": {
							Type:        schema.TypeString,
							Description: "Name of the storage.",
							Computed:    true,
							Default:    "AWS - S3 Bucket",
						},
						"s3_bucket_configuration": {
							Type:        schema.TypeSet,
							Description: "Configuration for S3 bucket.",
							Required:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket_name": {
										Type:        schema.TypeString,
										Description: "Name of the S3 bucket.",
										Required:    true,
									},
									"cloud_name": {
										Type:        schema.TypeString,
										Description: "Name of the cloud provider.",
										Computed:    true,
										Default:    "AWS",
									},
									"aws_region": {
										Type:        schema.TypeString,
										Description: "Region of the AWS S3 bucket.",
										Required:    true,
									},
								},
							},
						},
					},
				},
    		},
		},
	}
}

func resourceClassificationCreateContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*Client)

	// check provided fields against schema
	// if isOk, err := checkResourceRequiredFields(requiredSecretManagerFieldsJson, ignoreSecretManagerParamsByServerType, d); !isOk {
	// 	return diag.FromErr(err)
	// }

	// convert provided fields into API payload
	classification := ResourceWrapper{}
	classificationType := d.Get("type").(string)
	createIntegrationResource(&classification, classificationType, d)

	// create resource
	log.Printf("[INFO] Creating Classification of type: %s\n", classificationType)
	createClassificationResponse, err := client.CreateClassification(classification)
	if err != nil {
		log.Printf("[ERROR] adding classification of type: %s | err: %s\n", classificationType, err)
		return diag.FromErr(err)
	}

	// get asset_id
	// assetId := d.Get("asset_id").(string)

	// wait for remoteSyncState
	// err = waitForRemoteSyncState(ctx, dsfSecretManagerResourceType, assetId, m)
	// if err != nil {
	// 	diags = append(diags, diag.Diagnostic{
	// 		Severity: diag.Warning,
	// 		Summary:  fmt.Sprintf("Error while waiting for remoteSyncState = \"SYNCED\" for asset: %s", assetId),
	// 		Detail:   fmt.Sprintf("Error: %s\n", err),
	// 	})
	// }

	// set ID
	classificationId := createClassificationResponse.Data.ID
	d.SetId(classificationId)

	// Set the rest of the state from the resource read
	resourceClassificationReadContext(ctx, d, m)

	return diags
}

// TODO
func resourceClassificationReadContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	classificationId := d.Id()

	log.Printf("[INFO] Reading classification with Id: %s\n", classificationId)

	classificationReadResponse, err := client.ReadClassification(classificationId)

	if err != nil {
		log.Printf("[ERROR] Reading classificationReadResponse with classificationId: %s | err: %s\n", classificationId, err)
		return diag.FromErr(err)
	}

	if classificationReadResponse != nil {
		log.Printf("[INFO] Reading Classification with classificationId: %s | err: %s\n", classificationId, err)
	}

	log.Printf("[DEBUG] classificationReadResponse: %s\n", classificationReadResponse.Data.IntegrationData.ID)
	// Set returned and computed values
	d.Set("id", classificationReadResponse.Data.IntegrationData.ID)
	d.Set("description", classificationReadResponse.Data.IntegrationData.Description)
	d.Set("type", classificationReadResponse.Data.IntegrationData.Type)
	d.Set("status", classificationReadResponse.Data.IntegrationData.Status)
	d.Set("display_name", classificationReadResponse.Data.IntegrationData.DisplayName)
	d.Set("last_status_update", classificationReadResponse.Data.IntegrationData.LastStatusUpdate)
	// if secretManagerReadResponse.Data.AssetData.ServerPort != nil {
	// 	var serverPort string
	// 	if serverPortNum, ok := secretManagerReadResponse.Data.AssetData.ServerPort.(float64); ok {
	// 		serverPort = fmt.Sprintf("%d", int(serverPortNum))
	// 	} else {
	// 		serverPort = secretManagerReadResponse.Data.AssetData.ServerPort.(string)
	// 	}
	// 	d.Set("server_port", serverPort)
	// }
	// d.Set("server_type", secretManagerReadResponse.Data.ServerType)
	// d.Set("used_for", secretManagerReadResponse.Data.AssetData.UsedFor)
	// d.Set("version", secretManagerReadResponse.Data.AssetData.Version)

	// if secretManagerReadResponse.Data.AssetData.AwsProxyConfig != nil {
	// 	awsProxyConfig := &schema.Set{F: resourceAssetDataAWSProxyConfigHash}
	// 	awsProxyConfigMap := map[string]interface{}{}
	// 	awsProxyConfigMap["http"] = secretManagerReadResponse.Data.AssetData.AwsProxyConfig.HTTP
	// 	awsProxyConfigMap["https"] = secretManagerReadResponse.Data.AssetData.AwsProxyConfig.HTTPS
	// 	awsProxyConfig.Add(awsProxyConfigMap)
	// 	d.Set("aws_proxy_config", awsProxyConfig)
	// }

	// if secretManagerReadResponse.Data.AssetData.ServiceEndpoints != nil {
	// 	serviceEndpoints := &schema.Set{F: resourceAssetDataServiceEndpointsHash}
	// 	serviceEndpointsMap := map[string]interface{}{}
	// 	serviceEndpointsMap["logs"] = secretManagerReadResponse.Data.AssetData.ServiceEndpoints.Logs
	// 	serviceEndpoints.Add(serviceEndpointsMap)
	// 	d.Set("service_endpoints", serviceEndpoints)
	// }

	// connections := &schema.Set{F: resourceSecretManagerConnectionHash}
	// for _, v := range secretManagerReadResponse.Data.AssetData.Connections {
	// 	connection := map[string]interface{}{}
	// 	connection["access_id"] = v.ConnectionData.AccessID
	// 	connection["aws_iam_server_id"] = v.ConnectionData.AwsIamServerID
	// 	connection["ca_certs_path"] = v.ConnectionData.CaCertsPath
	// 	connection["cert_file"] = v.ConnectionData.CaFile
	// 	connection["credential_expiry"] = v.ConnectionData.CredentialExpiry
	// 	connection["external_id"] = v.ConnectionData.ExternalID
	// 	connection["key_file"] = v.ConnectionData.KeyFile
	// 	connection["nonce"] = v.ConnectionData.Nonce
	// 	connection["protocol"] = v.ConnectionData.Protocol
	// 	connection["query"] = v.ConnectionData.Query
	// 	connection["reason"] = v.Reason
	// 	connection["region"] = v.ConnectionData.Region
	// 	connection["role_name"] = v.ConnectionData.RoleName
	// 	connection["secret_key"] = v.ConnectionData.SecretKey
	// 	connection["self_signed"] = v.ConnectionData.SelfSigned
	// 	connection["ssl"] = v.ConnectionData.Ssl
	// 	connection["store_aws_credentials"] = v.ConnectionData.StoreAwsCredentials
	// 	connection["username"] = v.ConnectionData.Username
	// 	connection["v2_key_engine"] = v.ConnectionData.V2KeyEngine

	// 	// Handle structs
	// 	if v.ConnectionData.AmazonSecret != nil {
	// 		amazonSecret := &schema.Set{F: resourceConnectionDataAmazonSecretHash}
	// 		amazonSecretMap := map[string]interface{}{}
	// 		//amazonSecretMap["field_mapping"] = v.ConnectionData.AmazonSecret.FieldMapping
	// 		amazonSecretMap["secret_asset_id"] = v.ConnectionData.AmazonSecret.SecretAssetID
	// 		amazonSecretMap["secret_name"] = v.ConnectionData.AmazonSecret.SecretName
	// 		amazonSecret.Add(amazonSecretMap)
	// 		connection["amazon_secret"] = amazonSecret
	// 	}

	// 	if v.ConnectionData.CredentialFields != nil {
	// 		credentialFields := &schema.Set{F: resourceConnectionDataCredentialFieldsHash}
	// 		credentialFieldsMap := map[string]interface{}{}
	// 		credentialFieldsMap["credential_source"] = v.ConnectionData.CredentialFields.CredentialSource
	// 		credentialFieldsMap["role_arn"] = v.ConnectionData.CredentialFields.RoleArn
	// 		credentialFields.Add(credentialFieldsMap)
	// 		connection["credential_fields"] = credentialFields
	// 	}

	// 	if v.ConnectionData.CyberarkSecret != nil {
	// 		amazonSecret := &schema.Set{F: resourceConnectionDataCyberarkSecretHash}
	// 		amazonSecretMap := map[string]interface{}{}
	// 		//amazonSecretMap["field_mapping"] = v.ConnectionData.AmazonSecret.FieldMapping
	// 		amazonSecretMap["secret_asset_id"] = v.ConnectionData.CyberarkSecret.SecretAssetID
	// 		amazonSecretMap["secret_name"] = v.ConnectionData.CyberarkSecret.SecretName
	// 		amazonSecret.Add(amazonSecretMap)
	// 		connection["cyberark_secret"] = amazonSecret
	// 	}

	// 	if v.ConnectionData.HashicorpSecret != nil {
	// 		hashicorpSecret := &schema.Set{F: resourceConnectionDataHashicorpSecretHash}
	// 		hashicorpSecretMap := map[string]interface{}{}
	// 		//hashicorpSecretMap["field_mapping"] = v.ConnectionData.HashicorpSecret.Path
	// 		hashicorpSecretMap["path"] = v.ConnectionData.HashicorpSecret.Path
	// 		hashicorpSecretMap["secret_asset_id"] = v.ConnectionData.HashicorpSecret.SecretAssetID
	// 		hashicorpSecretMap["secret_name"] = v.ConnectionData.HashicorpSecret.SecretName
	// 		hashicorpSecret.Add(hashicorpSecretMap)
	// 		connection["hashicorp_secret"] = hashicorpSecret
	// 	}

	// 	connections.Add(connection)
	// }
	// d.Set("asset_connection", connections)

	log.Printf("[INFO] Finished reading classification with classificationId: %s\n", classificationId)

	return nil
}

func resourceClassificationUpdateContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*Client)

	// check provided fields against schema
	classificationId := d.Id()
	// if isOk, err := checkResourceRequiredFields(requiredSecretManagerFieldsJson, ignoreSecretManagerParamsByServerType, d); !isOk {
	// 	return diag.FromErr(err)
	// }

	// convert provided fields into API payload
	classification := ResourceWrapper{}
	classificationType := d.Get("type").(string)
	createIntegrationResource(&classification, classificationType, d)

	// update resource
	log.Printf("[INFO] Updating classification for Type: %s and Id: %s\n", classification.Data.IntegrationData.Type, classification.Data.IntegrationData.ID)
	_, err := client.UpdateClassification(classificationId, classification)
	if err != nil {
		log.Printf("[ERROR] Updating classification for Type: %s and Id: %s | err:%s\n", classification.Data.IntegrationData.Type, classification.Data.IntegrationData.ID, err)
		return diag.FromErr(err)
	}

	// get asset_id
	// assetId := d.Get("asset_id").(string)

	// wait for remoteSyncState
	// err = waitForRemoteSyncState(ctx, dsfSecretManagerResourceType, assetId, m)
	// if err != nil {
	// 	diags = append(diags, diag.Diagnostic{
	// 		Severity: diag.Warning,
	// 		Summary:  fmt.Sprintf("Error while waiting for remoteSyncState = \"SYNCED\" for asset: %s", assetId),
	// 		Detail:   fmt.Sprintf("Error: %s\n", err),
	// 	})
	// }

	// set ID
	d.SetId(classificationId)

	// Set the rest of the state from the resource read
	resourceClassificationReadContext(ctx, d, m)

	return diags
}

func resourceClassificationDeleteContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	classificationId := d.Id()

	log.Printf("[INFO] Deleting classification with classificationId: %s", classificationId)

	classificationDeleteResponse, err := client.DeleteClassification(classificationId)
	if classificationDeleteResponse != nil {
		log.Printf("[INFO] DSF classification has already been deleted with classificationId: %s | err: %s\n", classificationId, err)
	}

	return nil
}

// TODO
func resourceIntegrationDatabaseDetailsHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	// if v, ok := m["id"]; ok {
	// 	buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	// }

	// if v, ok := m["description"]; ok {
	// 	buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	// }

	// if v, ok := m["type"]; ok {
	// 	buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	// }

	// if v, ok := m["status"]; ok {
	// 	buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	// }

	// if v, ok := m["display_name"]; ok {
	// 	buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	// }

	// if v, ok := m["last_status_update"]; ok {
	// 	buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	// }

	if v, ok := m["database_details"]; ok {
		databaseDetails := v.(*schema.Set).List()
		for _, databaseDetail := range databaseDetails {
			databaseDetailMap := databaseDetail.(map[string]interface{})
			if v, ok := databaseDetailMap["database_type"]; ok {
				buf.WriteString(fmt.Sprintf("%s-", v.(string)))
			}
			if v, ok := databaseDetailMap["mongo_configuration"]; ok {
				buf.WriteString(fmt.Sprintf("%s-", v.(string)))
			}
		}
	}
	if v, ok := m["storage_details"]; ok {
		storageDetails := v.(*schema.Set).List()
		for _, storageDetail := range storageDetails {
			storageDetailMap := storageDetail.(map[string]interface{})
			if v, ok := storageDetailMap["storage_type"]; ok {
				buf.WriteString(fmt.Sprintf("%s-", v.(string)))
			}
			if v, ok := storageDetailMap["s3_bucket_configuration"]; ok {
				s3BucketDetails := v.(*schema.Set).List()
				for _, s3BucketDetail := range s3BucketDetails {
					s3BucketDetailMap := s3BucketDetail.(map[string]interface{})
					if v, ok := s3BucketDetailMap["bucket_name"]; ok {
						buf.WriteString(fmt.Sprintf("%s-", v.(string)))
					}
					if v, ok := s3BucketDetailMap["cloud_name"]; ok {
						buf.WriteString(fmt.Sprintf("%s-", v.(string)))
					}
					if v, ok := s3BucketDetailMap["aws_region"]; ok {
						buf.WriteString(fmt.Sprintf("%s-", v.(string)))
					}
				}
			}
		}
	}

	return PositiveHash(buf.String())
}
