package bthtml

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type payloadT struct {
	Cmd                    string  `json:"cmd"`
	URL                    string  `json:"url"`
	MaxTimeoutMilliseconds int     `json:"maxTimeout"`
	WaitInSeconds          float32 `json:"waitInSeconds"`
}

func GetFromURL(flareSolverrURL, webPageURL string, loadSeconds float32) (string, error) {
	client := &http.Client{}

	payload := payloadT{
		Cmd:                    "request.get",
		URL:                    webPageURL,
		MaxTimeoutMilliseconds: 60000,
		WaitInSeconds:          loadSeconds,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", errors.New("unable to marshal FlareSolverr payload to JSON: " + err.Error())
	}

	req, _ := http.NewRequest(
		"POST",
		flareSolverrURL,
		bytes.NewBuffer(payloadBytes),
	)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("unable to make FlareSolverr request to " + flareSolverrURL + ": " + err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("unable to read FlareSolverr response body: " + err.Error())
	}
	return string(body), nil
}
