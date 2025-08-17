package xenditpayment

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type XenditClient struct {
	APIKEY string
	Client *http.Client
}

func NewXenditClient(apiKey string) *XenditClient {
	return &XenditClient{
		APIKEY: apiKey,
		Client: &http.Client{},
	}
}

func (c *XenditClient) Post(path string, payload any) ([]byte, error) {
	url := os.Getenv("XENDIT_API_URL")

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.APIKEY, "")
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBodyClose, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return nil, err
	}

	return respBodyClose, nil
}
