package shortener

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ApiResponse struct {
	Status    int    `json:"status"`
	Message   string `json:"message"`
	Short_url string `json:"short_url"`
	Long_url  string `json:"long_url"`
}

type ApiErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func ShortenURL(long_url string, short_url string) (string, error) {
	url := fmt.Sprintf("https://1pt.one/shorten?short=%s&long=%s", short_url, long_url)

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("request failed: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errorResponse ApiErrorResponse
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)

		if err != nil {
			return "", fmt.Errorf("decoding response: %s", err.Error())
		}

		return "", fmt.Errorf("%s", errorResponse.Message)
	}

	var response ApiResponse
	err = json.NewDecoder(resp.Body).Decode(&response)

	if err != nil {
		return "", fmt.Errorf("decoding response: %s", err.Error())
	}

	return "1pt.one/" + response.Short_url, nil
}
