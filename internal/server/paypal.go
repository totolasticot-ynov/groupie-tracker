package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

const (
	paypalClientID = "Aao_IAK9WsSbSKqMd-HfOea_SwHvbJAaeJjpXC8eOmwNm5sj6s6kOLUoRSxOaTsnhR8Dr7oflFu2hj4e"
	paypalSecret   = "EGCbuKbWGfDNLQtzBx9BCBKj3RiMQ8p0AoBRCTtAKgNc1CHK8Y2xkjmnPAQPiUbXrRhuIphU7Hklj0Wg"
)

func getPayPalToken() string {
	req, _ := http.NewRequest(
		"POST",
		"https://api-m.sandbox.paypal.com/v1/oauth2/token",
		bytes.NewBufferString("grant_type=client_credentials"),
	)

	req.SetBasicAuth(paypalClientID, paypalSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
	}

	json.NewDecoder(resp.Body).Decode(&result)
	return result.AccessToken
}

func CreateOrder(w http.ResponseWriter, r *http.Request) {
	token := getPayPalToken()

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

func CaptureOrder(w http.ResponseWriter, r *http.Request) {
	var data struct {
		OrderID string `json:"orderID"`
	}

	json.NewDecoder(r.Body).Decode(&data)
	token := getPayPalToken()

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
