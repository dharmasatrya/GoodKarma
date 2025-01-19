package external

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/dharmasatrya/goodkarma/payment-service/entity"
)

func CreateXenditInvoice(req entity.XenditInvoiceRequest) (*entity.XenditInvoiceResponse, error) {
	// Prepare the request payload
	xenditUrl := os.Getenv("XENDIT_INVOICE_URL")
	payload := map[string]interface{}{
		"external_id": req.ExternalId,
		"amount":      req.Amount,
		"description": req.Description,
		"customer": map[string]interface{}{
			"given_names":   req.FirstName,
			"surname":       req.LastName,
			"email":         req.Email,
			"mobile_number": req.Phone,
		},
		"customer_notification_preference": map[string]interface{}{
			"invoice_created": []string{"whatsapp", "email"},
			"invoice_paid":    []string{"whatsapp", "email"},
		},
		"currency": "IDR",
	}

	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("37")
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Create the request
	request, err := http.NewRequest("POST", xenditUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("43")
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Get API key from environment variable and encode it
	apiKey := os.Getenv("XENDIT_API_KEY")
	if apiKey == "" {
		fmt.Println("49")
		return nil, fmt.Errorf("XENDIT_API_KEY not found in environment variables")
	}
	encodedKey := base64.StdEncoding.EncodeToString([]byte(apiKey))

	// Set headers
	request.Header.Set("Authorization", "Basic "+encodedKey)
	request.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Parse the response
	var result entity.XenditInvoiceResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("68")
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	// Check if response indicates an error
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		fmt.Println("74", result)
		return nil, fmt.Errorf("API request failed with status %d: %v", resp.StatusCode, result)
	}

	fmt.Println(&result)

	return &result, nil
}

func CreateXenditDisbursement(disbursementReq entity.XenditDisbursementRequest) (map[string]interface{}, error) {
	// Prepare the request payload
	xenditUrl := os.Getenv("XENDIT_DISBURSEMENT_URL")
	payload := map[string]interface{}{
		"external_id":         disbursementReq.ExternalId,
		"amount":              disbursementReq.Amount,
		"bank_code":           disbursementReq.BankCode,
		"account_holder_name": disbursementReq.AccountHolderName,
		"account_number":      disbursementReq.BankAccountNumber,
		"description":         disbursementReq.Description,
		"email_to":            []string{disbursementReq.Email},
	}

	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("error marshaling JSON: %v", err)
	}

	// Create the request
	request, err := http.NewRequest("POST", xenditUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Get API key from environment variable and encode it
	apiKey := os.Getenv("XENDIT_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("XENDIT_API_KEY not found in environment variables")
	}
	encodedKey := base64.StdEncoding.EncodeToString([]byte(apiKey))

	// Set headers
	request.Header.Set("Authorization", "Basic "+encodedKey)
	request.Header.Set("Content-Type", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Parse the response
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	// Check if response indicates an error
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API request failed with status %d: %v", resp.StatusCode, result)
	}

	return result, nil
}
