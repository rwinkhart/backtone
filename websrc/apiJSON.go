package websrc

import (
	"bytes"
	"errors"
	"io"
	"net/http"
)

type payloadFlareSolverrT struct {
	Cmd                    string  `json:"cmd"`
	URL                    string  `json:"url"`
	MaxTimeoutMilliseconds int     `json:"maxTimeout"`
	WaitInSeconds          float32 `json:"waitInSeconds"`
}

func GetByAPIJson(method, endpoint string, payload []byte) (string, error) {
	client := &http.Client{}

	req, _ := http.NewRequest(method, endpoint, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("unable to make " + method + " request to \"" + endpoint + "\": " + err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("unable to read \"" + endpoint + "\" response body: " + err.Error())
	}
	return string(body), nil
}
