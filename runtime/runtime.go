package runtime

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *ExtensionClient) Register(ctx context.Context, filename string) (*RegistrationResponse, error) {
	registrationUrl := fmt.Sprintf("%s/%s/extension/register", c.URL, c.Version)
	data := []RegistrationEvent{Invoke}
	jsonData, err := json.Marshal(data)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, registrationUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set(extensionNameHeader, filename)
	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	registrationResponse := RegistrationResponse{}
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(respBodyBytes, &registrationResponse)

	c.ExtensionId = resp.Header.Get(extensionIdentiferHeader)
	fmt.Println("Extension id:", c.ExtensionId)
	return &registrationResponse, nil
}

func (c *ExtensionClient) NextEvent(ctx context.Context) (*NextEventResponse, error) {
	nextUrl := fmt.Sprintf("%s/%s/extension/event/next", c.URL, c.Version)

	httpReq, err := http.NewRequestWithContext(ctx, "GET", nextUrl, nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set(extensionIdentiferHeader, c.ExtensionId)
	httpRes, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	if httpRes.StatusCode != 200 {
		return nil, fmt.Errorf("request failed with status %s", httpRes.Status)
	}
	defer httpRes.Body.Close()
	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return nil, err
	}
	res := NextEventResponse{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}
