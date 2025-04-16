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

// CreateCiphertrust adds an ciphertrust to DSF
func (c *Client) CreateCiphertrust(ciphertrust ResourceWrapper) (*ResourceWrapper, error) {
	log.Printf("[INFO] Adding  Ciphertrust Type: %s | ID: %s\n", ciphertrust.Data.IntegrationData.Type, ciphertrust.Data.IntegrationData.ID)

	//dsfDataSource := DSFDataSource{}
	ciphertrustJSON, err := json.Marshal(ciphertrust)
	log.Printf("[DEBUG] Adding ciphertrust - JSON: %s\n", ciphertrustJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to JSON marshal ciphertrust: %s\n", err)
	}

	resp, err := c.MakeCallWithQueryParams(http.MethodPost, endpointCiphertrust, "integration",ciphertrustJSON, c.config.Params)
	if err != nil {
		return nil, fmt.Errorf("error adding ciphertrust of type: %s | err: %s\n", ciphertrust.Data.IntegrationData.Type, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] Add DSF ciphertrust JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var createCiphertrustResponse ResourceWrapper
	err = json.Unmarshal([]byte(responseBody), &createCiphertrustResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing add Ciphertrust - JSON response type: %s | err: %s\n", ciphertrust.Data.IntegrationData.Type, err)
	}
	if createCiphertrustResponse.Errors != nil {
		return nil, fmt.Errorf("errors found in json response: %s", responseBody)
	}
	return &createCiphertrustResponse, nil
}

// ReadCiphertrust gets the DSF data source by ID
func (c *Client) ReadCiphertrust(ciphertrustId string) (*ResourceWrapper, error) {
	log.Printf("[INFO] Getting Ciphertrust for ciphertrustId: %s)\n", ciphertrustId)

	reqURL := fmt.Sprintf(endpointCiphertrust+"/%s", url.PathEscape(ciphertrustId))
	resp, err := c.MakeCall(http.MethodGet, reqURL, "integration",nil)
	if err != nil {
		return nil, fmt.Errorf("error reading Ciphertrust for ciphertrustId: %s | err: %s\n", ciphertrustId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] DSF Ciphertrust JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var readCiphertrustResponse ResourceWrapper
	err = json.Unmarshal([]byte(responseBody), &readCiphertrustResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing Ciphertrust JSON response for ciphertrustId: %s | Ciphertrust: %s err: %s\n", ciphertrustId, responseBody, err)
	}
	if readCiphertrustResponse.Errors != nil {
		return nil, fmt.Errorf("errors found in json response: %s", responseBody)
	}

	return &readCiphertrustResponse, nil
}

// ReadCiphertrusts gets all Ciphertrusts
func (c *Client) ReadCiphertrusts() (*ResourcesWrapper, error) {
	log.Printf("[INFO] Getting Ciphertrusts\n")

	resp, err := c.MakeCall(http.MethodGet, "get-config", "integration",nil)
	if err != nil {
		return nil, fmt.Errorf("error reading Ciphertrusts | err: %s\n", err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] DSF Ciphertrusts JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var readCiphertrustsResponse ResourcesWrapper
	err = json.Unmarshal([]byte(responseBody), &readCiphertrustsResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing Ciphertrusts JSON response: %s err: %s\n", responseBody, err)
	}
	if readCiphertrustsResponse.Errors != nil {
		return nil, fmt.Errorf("errors found in json response: %s", responseBody)
	}

	return &readCiphertrustsResponse, nil
}

// UpdateCiphertrust will update a specific ciphertrust record in DSF referenced by the dataSourceId
func (c *Client) UpdateCiphertrust(ciphertrustId string, ciphertrust ResourceWrapper) (*ResourceWrapper, error) {
	log.Printf("[INFO] Updating Ciphertrust with ciphertrustId: %s)\n", ciphertrustId)

	ciphertrustJSON, err := json.Marshal(ciphertrust)
	log.Printf("[DEBUG] Adding Ciphertrust - JSON: %s\n", ciphertrustJSON)
	if err != nil {
		return nil, fmt.Errorf("failed to JSON marshal Ciphertrust: %s", err)
	}

	reqURL := fmt.Sprintf(endpointCiphertrust+"/%s", url.PathEscape(ciphertrustId))
	resp, err := c.MakeCallWithQueryParams(http.MethodPut, reqURL, "integration",ciphertrustJSON, c.config.Params)
	if err != nil {
		return nil, fmt.Errorf("error updating Ciphertrust with ciphertrustId: %s | err: %s\n", ciphertrustId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] DSF update Ciphertrust JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var updateCiphertrustResponse ResourceWrapper
	err = json.Unmarshal([]byte(responseBody), &updateCiphertrustResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing update Ciphertrust JSON response for ciphertrustId: %s | err: %s\n", ciphertrustId, err)
	}
	if updateCiphertrustResponse.Errors != nil {
		return nil, fmt.Errorf("errors found in json response: %s", responseBody)
	}

	return &updateCiphertrustResponse, nil
}

// DeleteCiphertrust deletes an ciphertrust in DSF
func (c *Client) DeleteCiphertrust(ciphertrustId string) (*ResourceResponse, error) {
	log.Printf("[INFO] Deleting Ciphertrust with ciphertrustId: %s\n", ciphertrustId)

	reqURL := fmt.Sprintf(endpointCiphertrust+"/%s", url.PathEscape(ciphertrustId))
	resp, err := c.MakeCall(http.MethodDelete, reqURL,"integration", nil)
	if err != nil {
		return nil, fmt.Errorf("error deleting Ciphertrust for ciphertrustId: %s, %s\n", ciphertrustId, err)
	}

	// Read the body
	defer resp.Body.Close()
	responseBody, err := ioutil.ReadAll(resp.Body)

	// Dump JSON
	log.Printf("[DEBUG] DSF delete Ciphertrust with JSON response: %s\n", string(responseBody))

	// Parse the JSON
	var deleteCiphertrustResponse ResourceResponse
	err = json.Unmarshal([]byte(responseBody), &deleteCiphertrustResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing delete Ciphertrust JSON response for ciphertrustId: %s, %s\n", ciphertrustId, err)
	}
	if deleteCiphertrustResponse.Errors != nil {
		return nil, fmt.Errorf("errors found in json response: %s", responseBody)
	}

	return &deleteCiphertrustResponse, nil
}
