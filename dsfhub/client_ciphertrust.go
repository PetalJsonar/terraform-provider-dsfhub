package dsfhub

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

const endpointCiphertrust = "/ciphertrust"

// CreateCiphertrust adds a CipherTrust integration to DSF
func (c *Client) CreateCiphertrust(ciphertrust IntegrationResourceWrapper) (*IntegrationResourceWrapper, error) {
	log.Printf("[INFO] Adding Ciphertrust Type: %s | ID: %s\n", ciphertrust.IntegrationData.Type, ciphertrust.IntegrationData.ID)

	ciphertrustJSON, err := json.Marshal(ciphertrust)
	log.Printf("[DEBUG] Adding ciphertrust - JSON: %s\n", ciphertrustJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to JSON marshal ciphertrust: %s\n", err)
	}

	resp, err := c.MakeCall(http.MethodPost, endpointCiphertrust, "integration", ciphertrustJSON)
	if err != nil {
		return nil, fmt.Errorf("error adding ciphertrust of type: %s | err: %s\n", ciphertrust.IntegrationData.Type, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Add DSF ciphertrust JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var createCiphertrustResponse IntegrationResourceWrapper
	err = json.Unmarshal([]byte(responseBody), &createCiphertrustResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing add Ciphertrust - JSON response type: %s | err: %s\n", ciphertrust.IntegrationData.Type, err)
	}
	if createCiphertrustResponse.Error.Code != 200 {
		return nil, fmt.Errorf("errors found in json response: %s", responseBody)
	}
	return &createCiphertrustResponse, nil
}

// ReadCiphertrust gets the CipherTrust integration by ID
func (c *Client) ReadCiphertrust(ciphertrustId string) (*IntegrationResourceWrapper, error) {
	log.Printf("[INFO] Getting Ciphertrust for ciphertrustId: %s)\n", ciphertrustId)

	reqURL := fmt.Sprintf(endpointCiphertrust+"/%s", url.PathEscape(ciphertrustId))
	resp, err := c.MakeCall(http.MethodGet, reqURL, "integration", nil)
	if err != nil {
		return nil, fmt.Errorf("error reading Ciphertrust for ciphertrustId: %s | err: %s\n", ciphertrustId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] DSF Ciphertrust JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var readCiphertrustResponse IntegrationResourceWrapper
	err = json.Unmarshal([]byte(responseBody), &readCiphertrustResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing Ciphertrust JSON response for ciphertrustId: %s | Ciphertrust: %s err: %s\n", ciphertrustId, responseBody, err)
	}
	if readCiphertrustResponse.Error.Code != 200 {
		return nil, fmt.Errorf("errors found in json response: %s", responseBody)
	}

	return &readCiphertrustResponse, nil
}

// ReadCiphertrusts gets all CipherTrust integrations
func (c *Client) ReadCiphertrusts() (*IntegrationResourcesWrapper, error) {
	log.Printf("[INFO] Getting Ciphertrusts\n")

	resp, err := c.MakeCall(http.MethodGet, "get-config", "integration", nil)
	if err != nil {
		return nil, fmt.Errorf("error reading Ciphertrusts | err: %s\n", err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] DSF Ciphertrusts JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var readCiphertrustsResponse IntegrationResourcesWrapper
	err = json.Unmarshal([]byte(responseBody), &readCiphertrustsResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing Ciphertrusts JSON response: %s err: %s\n", responseBody, err)
	}
	if readCiphertrustsResponse.Error.Code != 200 {
		return nil, fmt.Errorf("errors found in json response: %s", responseBody)
	}

	return &readCiphertrustsResponse, nil
}

// UpdateCiphertrust will update a specific CipherTrust integration in DSF referenced by the ciphertrustId
func (c *Client) UpdateCiphertrust(ciphertrustId string, ciphertrust IntegrationResourceWrapper) (*IntegrationResourceWrapper, error) {
	log.Printf("[INFO] Updating Ciphertrust with ciphertrustId: %s)\n", ciphertrustId)

	ciphertrustJSON, err := json.Marshal(ciphertrust)
	log.Printf("[DEBUG] Adding Ciphertrust - JSON: %s\n", ciphertrustJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to JSON marshal Ciphertrust: %s", err)
	}

	reqURL := fmt.Sprintf(endpointCiphertrust+"/%s", url.PathEscape(ciphertrustId))
	resp, err := c.MakeCall(http.MethodPut, reqURL, "integration", ciphertrustJSON)
	if err != nil {
		return nil, fmt.Errorf("error updating Ciphertrust with ciphertrustId: %s | err: %s\n", ciphertrustId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] DSF update Ciphertrust JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var updateCiphertrustResponse IntegrationResourceWrapper
	err = json.Unmarshal([]byte(responseBody), &updateCiphertrustResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing update Ciphertrust JSON response for ciphertrustId: %s | err: %s\n", ciphertrustId, err)
	}
	if updateCiphertrustResponse.Error.Code != 200 {
		return nil, fmt.Errorf("errors found in json response: %s", responseBody)
	}

	return &updateCiphertrustResponse, nil
}

// DeleteCiphertrust deletes a CipherTrust integration in DSF
func (c *Client) DeleteCiphertrust(ciphertrustId string) (*IntegrationResourceWrapper, error) {
	log.Printf("[INFO] Deleting Ciphertrust with ciphertrustId: %s\n", ciphertrustId)

	reqURL := fmt.Sprintf(endpointCiphertrust+"/%s", url.PathEscape(ciphertrustId))

	relevantParams := map[string]interface{}{
		"acknowledgeDeletionImpact": c.config.Params["acknowledgeDeletionImpact"],
	}

	resp, err := c.MakeCallWithQueryParams(http.MethodDelete, reqURL, "integration", nil, relevantParams)
	if err != nil {
		return nil, fmt.Errorf("error deleting Ciphertrust for ciphertrustId: %s, %s\n", ciphertrustId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] DSF delete Ciphertrust with JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var deleteCiphertrustResponse IntegrationResourceWrapper // IntegrationError ?
	err = json.Unmarshal([]byte(responseBody), &deleteCiphertrustResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing delete Ciphertrust JSON response for ciphertrustId: %s, %s\n", ciphertrustId, err)
	}
	if deleteCiphertrustResponse.Error.Code != 200 {
		return nil, fmt.Errorf("errors found in json response: %s", responseBody)
	}

	return &deleteCiphertrustResponse, nil
}
