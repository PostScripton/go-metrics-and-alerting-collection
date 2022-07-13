package client

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Client struct {
	baseURI string
	client  http.Client
}

func New(baseURI string) Client {
	return Client{
		baseURI: baseURI,
		client: http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) UpdateMetric(metricType string, name string, value string) {
	fmt.Printf("--- [%s] \"%s\": %s\n", metricType, name, value)

	url := fmt.Sprintf("%s/update/%s/%s/%s", c.baseURI, metricType, name, value)
	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	request.Header.Set("Content-Type", "text/plain")

	response, err := c.client.Do(request)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if response.StatusCode != http.StatusOK {
		defer response.Body.Close()
		body, errBodyReader := io.ReadAll(response.Body)
		if errBodyReader != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Printf("Status code: %d\n", response.StatusCode)
		fmt.Printf("Message: %s\n", string(body))
		return
	}
}
