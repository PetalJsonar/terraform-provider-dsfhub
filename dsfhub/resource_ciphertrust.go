package dsfhub

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCiphertrust() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCiphertrustCreateContext,
		ReadContext:   resourceCiphertrustReadContext,
		UpdateContext: resourceCiphertrustUpdateContext,
		DeleteContext: resourceCiphertrustDeleteContext,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Description: "ID of the ciphertrust.",
				Optional:    true,
				Computed:    true,
				// Default: nil,
			},
			"description": {
				Type:        schema.TypeString,
				Description: "Description of the ciphertrust.",
				Optional:    true,
				Computed:    true,
				Default: "Used for integrating with Thales CipherTrust Manager capabilities.",
			},
			"type": {
				Type:        schema.TypeString,
				Description: "Type of the ciphertrust.",
				Computed:    true,
				Default: "CipherTrust Manager",
			},
			"status": {
				Type:        schema.TypeString,
				Description: "Status of the ciphertrust.",
				Optional:    true,
				Default: "N/A",
			},
			"display_name": {
				Type:        schema.TypeString,
				Description: "User-friendly name of the ciphertrust, defined by user.",
				Optional:    true,
			},
			"last_status_update": {
				Type:        schema.TypeString,
				Description: "Timestamp of the last status update.",
				Optional:    true,
				Default: nil,
			},
			"hostname": {
				Type:        schema.TypeString,
				Description: "Hostname of the ciphertrust.",
				Required:    true,
			},
			"port": {
				Type:        schema.TypeInt,
				Description: "Port of the ciphertrust.",
				Required:    true,
			},
			"username": {
				Type:        schema.TypeString,
				Description: "Username for authentication.",
				Optional:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "Password for authentication.",
				Optional:    true,
				Sensitive: true,
			},
			"cm_name": {
				Type:        schema.TypeString,
				Description: "Name of the Ciphertrust Manager.",
				Required:    true,
			},
			"is_load_balancer": {
				Type:        schema.TypeBool,
				Description: "Indicates if the ciphertrust is a load balancer.",
				Optional:    true,
				Default:     false,
			},
			"auth_method": {
				Type:        schema.TypeString,
				Description: "Authentication method for the ciphertrust.",
				Required:   true,
			},
			"registration_token": {
				Type:        schema.TypeString,
				Description: "Registration token for the ciphertrust.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func resourceCiphertrustCreateContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*Client)

	// check provided fields against schema
	// if isOk, err := checkResourceRequiredFields(requiredSecretManagerFieldsJson, ignoreSecretManagerParamsByServerType, d); !isOk {
	// 	return diag.FromErr(err)
	// }

	// convert provided fields into API payload
	ciphertrust := ResourceWrapper{}
	ciphertrustType := d.Get("type").(string)
	createIntegrationResource(&ciphertrust, ciphertrustType, d)

	// create resource
	log.Printf("[INFO] Creating Ciphertrust of type: %s\n", ciphertrustType)
	createCiphertrustResponse, err := client.CreateCiphertrust(ciphertrust)
	if err != nil {
		log.Printf("[ERROR] adding ciphertrust of type: %s | err: %s\n", ciphertrustType, err)
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
	ciphertrustId := createCiphertrustResponse.Data.ID
	d.SetId(ciphertrustId)

	// Set the rest of the state from the resource read
	resourceCiphertrustReadContext(ctx, d, m)

	return diags
}

// TODO
func resourceCiphertrustReadContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	ciphertrustId := d.Id()

	log.Printf("[INFO] Reading ciphertrust with Id: %s\n", ciphertrustId)

	ciphertrustReadResponse, err := client.ReadCiphertrust(ciphertrustId)

	if err != nil {
		log.Printf("[ERROR] Reading ciphertrustReadResponse with ciphertrustId: %s | err: %s\n", ciphertrustId, err)
		return diag.FromErr(err)
	}

	if ciphertrustReadResponse != nil {
		log.Printf("[INFO] Reading Ciphertrust with ciphertrustId: %s | err: %s\n", ciphertrustId, err)
	}

	log.Printf("[DEBUG] ciphertrustReadResponse: %s\n", ciphertrustReadResponse.Data.IntegrationData.ID)
	// Set returned and computed values
	d.Set("id", ciphertrustReadResponse.Data.IntegrationData.ID)
	d.Set("description", ciphertrustReadResponse.Data.IntegrationData.Description)
	d.Set("type", ciphertrustReadResponse.Data.IntegrationData.Type)
	d.Set("status", ciphertrustReadResponse.Data.IntegrationData.Status)
	d.Set("display_name", ciphertrustReadResponse.Data.IntegrationData.DisplayName)
	d.Set("last_status_update", ciphertrustReadResponse.Data.IntegrationData.LastStatusUpdate)
	d.Set("hostname", ciphertrustReadResponse.Data.IntegrationData.Hostname)
	d.Set("port", ciphertrustReadResponse.Data.IntegrationData.Port)
	d.Set("username", ciphertrustReadResponse.Data.IntegrationData.Username)
	d.Set("password", ciphertrustReadResponse.Data.IntegrationData.Password)
	d.Set("cm_name", ciphertrustReadResponse.Data.IntegrationData.CMName)
	d.Set("is_load_balancer", ciphertrustReadResponse.Data.IntegrationData.IsLoadBalancer)
	d.Set("auth_method", ciphertrustReadResponse.Data.IntegrationData.AuthMethod)
	d.Set("registration_token", ciphertrustReadResponse.Data.IntegrationData.RegistrationToken)

	log.Printf("[INFO] Finished reading ciphertrust with ciphertrustId: %s\n", ciphertrustId)

	return nil
}

func resourceCiphertrustUpdateContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	client := m.(*Client)

	// check provided fields against schema
	ciphertrustId := d.Id()
	// if isOk, err := checkResourceRequiredFields(requiredSecretManagerFieldsJson, ignoreSecretManagerParamsByServerType, d); !isOk {
	// 	return diag.FromErr(err)
	// }

	// convert provided fields into API payload
	ciphertrust := ResourceWrapper{}
	ciphertrustType := d.Get("type").(string)
	createIntegrationResource(&ciphertrust, ciphertrustType, d)

	// update resource
	log.Printf("[INFO] Updating ciphertrust for Type: %s and Id: %s\n", ciphertrust.Data.ServerType, ciphertrust.Data.IntegrationData.ID)
	_, err := client.UpdateCiphertrust(ciphertrustId, ciphertrust)
	if err != nil {
		log.Printf("[ERROR] Updating Ciphertrust for Type: %s and Id: %s | err:%s\n", ciphertrust.Data.ServerType, ciphertrust.Data.IntegrationData.ID, err)
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
	d.SetId(ciphertrustId)

	// Set the rest of the state from the resource read
	resourceCiphertrustReadContext(ctx, d, m)

	return diags
}

func resourceCiphertrustDeleteContext(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)
	ciphertrustId := d.Id()

	log.Printf("[INFO] Deleting ciphertrust with ciphertrustId: %s", ciphertrustId)

	ciphertrustDeleteResponse, err := client.DeleteCiphertrust(ciphertrustId)
	if ciphertrustDeleteResponse != nil {
		log.Printf("[INFO] DSF ciphertrust has already been deleted with ciphertrustId: %s | err: %s\n", ciphertrustId, err)
	}

	return nil
}

// TODO
func resourceCiphertrustDatabaseDetailsHash(v interface{}) int {
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
