package dsfhub

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceClassification() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClassificationRead,
		Description: "Provides Classification from a unique id.",

		Schema: map[string]*schema.Schema{
			// Computed Attributes
			"id": {
				Type:        schema.TypeString,
				Description: "Current ID",
				Required:    true,
				Optional:    false,
			},
		},
	}
}

func dataSourceClassificationRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	curIntegrationId := d.Get("id").(string)
	log.Printf("[INFO] DataSource - Reading Classification with integrationId: %s", curIntegrationId)

	classificationReadResponse, err := client.ReadClassification(curIntegrationId)
	if classificationReadResponse != nil {
		log.Printf("[INFO] Reading Classification with integrationId: %s | err: %s\n", curIntegrationId, err)
	}
	integrationId := classificationReadResponse.IntegrationData.ID
	d.Set("asset_id", integrationId)
	d.SetId(integrationId)

	log.Printf("[INFO] Finished reading DataSource Classification with integrationId: %s\n", integrationId)
	return nil
}

// func dataSourceClassifications() *schema.Resource {
// 	return &schema.Resource{
// 		ReadContext: dataSourceSecretManagersRead,
// 		Description: "Provides a list of SecretManagers filtering for asset_id values by regex.",

// 		Schema: map[string]*schema.Schema{
// 			// Computed Attributes
// 			"asset_id_regex": {
// 				Type:        schema.TypeString,
// 				Description: "Regex pattern for asset IDs",
// 				Optional:    true,
// 				Default:     nil,
// 			},
// 			"asset_ids": {
// 				Type:        schema.TypeList,
// 				Description: "List of asset IDs",
// 				Elem: &schema.Schema{
// 					Type: schema.TypeString,
// 				},
// 				Optional: false,
// 				Computed: true,
// 				Default:  nil,
// 			},
// 		},
// 	}
// }

// func dataSourceClassificationsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
// 	client := m.(*Client)

// 	assetIdRegex := d.Get("asset_id_regex").(string)
// 	log.Printf("[INFO] Data Source - Reading SecretManagers filtering for asset_ids with assetIdRegex: %s", assetIdRegex)

// 	secretManagersReadResponse, err := client.ReadClassifications()
// 	if secretManagersReadResponse != nil {
// 		log.Printf("[INFO] Data Source - Reading SecretManagers filtering for asset_ids with assetIdRegex: %s | err: %s\n", assetIdRegex, err)
// 	}
// 	var assetIds []string
// 	for _, ca := range secretManagersReadResponse.Data {
// 		match, _ := regexp.MatchString(assetIdRegex, ca.ID)
// 		log.Printf("[INFO] Checking asset_id: %v against regex: %v match:%v\n", ca.ID, assetIdRegex, match)
// 		if match {
// 			log.Printf("[INFO] Matched asset_id: %v against regex: %v match:%v\n", ca.ID, assetIdRegex, match)
// 			assetIds = append(assetIds, ca.ID)
// 		} else {
// 			log.Printf("[INFO] Did not match asset_id: %v against regex: %v match:%v\n", ca.ID, assetIdRegex, match)
// 		}
// 	}

// 	d.Set("asset_ids", assetIds)
// 	d.SetId(client.config.DSFHUBHost + "_secretManagers")

// 	log.Printf("[INFO] Finished reading Data Source SecretManagers with assetIdRegex: %s\n", assetIdRegex)
// 	return nil
// }
