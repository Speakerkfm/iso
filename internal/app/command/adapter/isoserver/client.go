package isoserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Speakerkfm/iso/internal/pkg/config"
	"github.com/Speakerkfm/iso/internal/pkg/models"
)

type Client struct {
}

// New ...
func New() *Client {
	return &Client{}
}

// GetServiceConfigs ...
func (c *Client) GetServiceConfigs(ctx context.Context) ([]models.ServiceConfigDesc, error) {
	resp, err := http.Get(fmt.Sprintf("http://%s/service_configs", config.ISOServerAdminHost))
	if err != nil {
		return nil, fmt.Errorf("fail to get service configs: %w", err)
	}
	defer resp.Body.Close()

	var res []models.ServiceConfigDesc
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("fail to decode service configs response: %w", err)
	}

	return res, nil
}

// SaveServiceConfigs ...
func (c *Client) SaveServiceConfigs(ctx context.Context, serviceConfigs []models.ServiceConfigDesc) error {
	reqBody := bytes.NewBuffer(nil)
	if err := json.NewEncoder(reqBody).Encode(serviceConfigs); err != nil {
		return fmt.Errorf("fail to encode service configs: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("http://%s/service_configs", config.ISOServerAdminHost), reqBody)
	if err != nil {
		return fmt.Errorf("fail to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("fail to send req to iso server: %w", err)
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	if resp.Body == nil {
		return fmt.Errorf("got bad status code from server: %d", resp.StatusCode)
	}

	respErr, _ := ioutil.ReadAll(resp.Body)
	return fmt.Errorf("got bad status code from server: %d, err: %s", resp.StatusCode, string(respErr))
}
