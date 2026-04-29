package services

import (
	"fmt"
	"io"
	"net/http"
)

func GraphApiGetRequest(url string, accessToken string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/json")
	// Utilisé par l'import d'évènements Outlook pour que le corps/la description de l'évènement soit en text plutôt qu'en HTML
	req.Header.Set("Prefer", "outlook.body-content-type=\"text\"")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("graph API error: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}
