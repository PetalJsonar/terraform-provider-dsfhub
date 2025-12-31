package dsfhub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const endpointClassification = "/classification"

// CreateClassification adds a classification integration to DSF
func (c *Client) CreateClassification(classification IntegrationResourceWrapper) (*IntegrationResourceWrapper, error) {
	log.Printf("[INFO] Adding Classification Type: %s | ID: %s\n", classification.IntegrationData.Type, classification.IntegrationData.ID)

	classificationJSON, err := json.Marshal(classification)
	log.Printf("[DEBUG] Adding classification - JSON: %s\n", classificationJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to JSON marshal classification: %s\n", err)
	}

	resp, err := c.MakeCall(http.MethodPost, endpointClassification, "integration", classificationJSON)
	if err != nil {
		return nil, fmt.Errorf("error adding classification of type: %s | err: %s\n", classification.IntegrationData.Type, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Add DSF classification JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var createClassificationResponse IntegrationResourceWrapper
	err = json.Unmarshal([]byte(responseBody), &createClassificationResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing add Classification - JSON response type: %s | err: %s\n", classification.IntegrationData.Type, err)
	}
	if createClassificationResponse.Error.Code != 200 {
		return nil, fmt.Errorf("errors found in json response: %s", responseBody)
	}
	return &createClassificationResponse, nil
}

// ReadClassification gets the classification integration by ID
func (c *Client) ReadClassification(classificationId string) (*IntegrationResourceWrapper, error) {
	log.Printf("[INFO] Getting Classification for classificationId: %s)\n", classificationId)

	reqURL := fmt.Sprintf(endpointClassification+"/%s", url.PathEscape(classificationId))
	resp, err := c.MakeCall(http.MethodGet, reqURL, "integration", nil)
	if err != nil {
		return nil, fmt.Errorf("error reading Classification for classificationId: %s | err: %s\n", classificationId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] DSF Classification JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var readClassificationResponse IntegrationResourceWrapper
	err = json.Unmarshal([]byte(responseBody), &readClassificationResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing Classification JSON response for classificationId: %s | Classification: %s err: %s\n", classificationId, responseBody, err)
	}
	if readClassificationResponse.Error.Code != 200 {
		return nil, fmt.Errorf("errors found in json response: %s", responseBody)
	}

	return &readClassificationResponse, nil
}

// ReadClassifications gets all Classification integrations
func (c *Client) ReadClassifications() (*IntegrationResourcesWrapper, error) {
	log.Printf("[INFO] Getting Classifications\n")

	resp, err := c.MakeCall(http.MethodGet, "get-config", "integration", nil)
	if err != nil {
		return nil, fmt.Errorf("error reading Classifications | err: %s\n", err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] DSF Classifications JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var readClassificationsResponse IntegrationResourcesWrapper
	err = json.Unmarshal([]byte(responseBody), &readClassificationsResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing Classifications JSON response: %s err: %s\n", responseBody, err)
	}
	if readClassificationsResponse.Error.Code != 200 {
		return nil, fmt.Errorf("errors found in json response: %s", responseBody)
	}

	return &readClassificationsResponse, nil
}

// UpdateClassification will update a specific classification record in DSF referenced by the classificationId
func (c *Client) UpdateClassification(classificationId string, classification IntegrationResourceWrapper) (*IntegrationResourceWrapper, error) {
	log.Printf("[INFO] Updating Classification with classificationId: %s)\n", classificationId)

	classificationJSON, err := json.Marshal(classification)
	log.Printf("[DEBUG] Adding Classification - JSON: %s\n", classificationJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to JSON marshal Classification: %s", err)
	}

	reqURL := fmt.Sprintf(endpointClassification+"/%s", url.PathEscape(classificationId))
	resp, err := c.MakeCall(http.MethodPut, reqURL, "integration", classificationJSON)
	if err != nil {
		return nil, fmt.Errorf("error updating Classification with classificationId: %s | err: %s\n", classificationId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] DSF update Classification JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var updateClassificationResponse IntegrationResourceWrapper
	err = json.Unmarshal([]byte(responseBody), &updateClassificationResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing update Classification JSON response for classificationId: %s | err: %s\n", classificationId, err)
	}
	if updateClassificationResponse.Error.Code != 200 {
		return nil, fmt.Errorf("errors found in json response: %s", responseBody)
	}

	return &updateClassificationResponse, nil
}

// DeleteClassification delete a classification integration in DSF
func (c *Client) DeleteClassification(classificationId string) (*IntegrationResourceWrapper, error) {
	log.Printf("[INFO] Deleting Classification with classificationId: %s\n", classificationId)

	reqURL := fmt.Sprintf(endpointClassification+"/%s", url.PathEscape(classificationId))

	relevantParams := map[string]interface{}{
		"forceDelete": c.config.Params["forceDelete"],
	}

	resp, err := c.MakeCallWithQueryParams(http.MethodDelete, reqURL, "integration", nil, relevantParams)
	if err != nil {
		return nil, fmt.Errorf("error deleting Classification for classificationId: %s, %s\n", classificationId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] DSF delete Classification with JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var deleteClassificationResponse IntegrationResourceWrapper // IntegrationError ?
	err = json.Unmarshal([]byte(responseBody), &deleteClassificationResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing delete Classification JSON response for classificationId: %s, %s\n", classificationId, err)
	}
	if deleteClassificationResponse.Error.Code != 200 {
		return nil, fmt.Errorf("errors found in json response: %s", responseBody)
	}

	return &deleteClassificationResponse, nil
}
