package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

const (
	paypalClientID = "TON_CLIENT_ID_SANDBOX_GROUPIE_TRACKER"
	paypalSecret   = "TON_SECRET_SANDBOX_GROUPIE_TRACKER"
)

func getPayPalToken() (string, error) {
	req, _ := http.NewRequest(
		"POST",
		"https://api-m.sandbox.paypal.com/v1/oauth2/token",
		bytes.NewBufferString("grant_type=client_credentials"),
	)

	req.SetBasicAuth(paypalClientID, paypalSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	return result.AccessToken, nil
}

func createOrderHandler(w http.ResponseWriter, r *http.Request) {
	token, err := getPayPalToken()
	if err != nil {
		http.Error(w, "PayPal token error", 500)
		return
	}

	body := `{
		"intent": "CAPTURE",
		"purchase_units": [{
			"amount": {
				"currency_code": "EUR",
				"value": "5.00"
			}
		}]
	}`

	req, _ := http.NewRequest(
		"POST",
		"https://api-m.sandbox.paypal.com/v2/checkout/orders",
		bytes.NewBufferString(body),
	)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	io.Copy(w, resp.Body)
}

func captureOrderHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		OrderID string `json:"orderID"`
	}
	json.NewDecoder(r.Body).Decode(&data)

	token, err := getPayPalToken()
	if err != nil {
		http.Error(w, "PayPal token error", 500)
		return
	}

	req, _ := http.NewRequest(
		"POST",
		"https://api-m.sandbox.paypal.com/v2/checkout/orders/"+data.OrderID+"/capture",
		nil,
	)

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	io.Copy(w, resp.Body)
}
