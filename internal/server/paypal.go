package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// --- Récupération des identifiants depuis les variables d'environnement ---
var (
	paypalClientID = "Aao_IAK9WsSbSKqMd-HfOea_SwHvbJAaeJjpXC8eOmwNm5sj6s6kOLUoRSxOaTsnhR8Dr7oflFu2hj4e" // Exemple: "Aao_IAK9WsSbSKqMd..."
	paypalSecret   = "EGCbuKbWGfDNLQtzBx9BCBKj3RiMQ8p0AoBRCTtAKgNc1CHK8Y2xkjmnPAQPiUbXrRhuIphU7Hklj0Wg" // Exemple: "TON_SECRET_SANDBOX"
)

// getPayPalToken récupère un token OAuth2 auprès de PayPal
func getPayPalToken() (string, error) {
	req, err := http.NewRequest(
		"POST",
		"https://api-m.sandbox.paypal.com/v1/oauth2/token",
		bytes.NewBufferString("grant_type=client_credentials"),
	)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(paypalClientID, paypalSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return "", &PayPalError{StatusCode: resp.StatusCode, Body: string(body)}
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.AccessToken, nil
}

// --- Handlers PayPal ---

// CreateOrderHandler crée un ordre PayPal et renvoie son ID
func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	token, err := getPayPalToken()
	if err != nil {
		http.Error(w, "Erreur token PayPal: "+err.Error(), http.StatusInternalServerError)
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

	req, err := http.NewRequest(
		"POST",
		"https://api-m.sandbox.paypal.com/v2/checkout/orders",
		bytes.NewBufferString(body),
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, string(body), resp.StatusCode)
		return
	}

	io.Copy(w, resp.Body)
}

// CaptureOrderHandler capture le paiement d’un ordre PayPal
func CaptureOrderHandler(w http.ResponseWriter, r *http.Request) {
	var data struct {
		OrderID string `json:"orderID"`
	}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, "Erreur lecture body: "+err.Error(), http.StatusBadRequest)
		return
	}

	token, err := getPayPalToken()
	if err != nil {
		http.Error(w, "Erreur token PayPal: "+err.Error(), http.StatusInternalServerError)
		return
	}

	req, err := http.NewRequest(
		"POST",
		"https://api-m.sandbox.paypal.com/v2/checkout/orders/"+data.OrderID+"/capture",
		nil,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		http.Error(w, string(body), resp.StatusCode)
		return
	}

	io.Copy(w, resp.Body)
}

// --- Types d’erreurs custom ---
type PayPalError struct {
	StatusCode int
	Body       string
}

func (e *PayPalError) Error() string {
	return e.Body
}
