package websrc

import (
	"encoding/json"
	"errors"
)

func GetByScrape(flareSolverrEndpoint, webPageURL string, loadSeconds float32) (string, error) {
	payload, err := json.Marshal(payloadFlareSolverrT{
		Cmd:                    "request.get",
		URL:                    webPageURL,
		MaxTimeoutMilliseconds: 60000,
		WaitInSeconds:          loadSeconds,
	})
	if err != nil {
		return "", errors.New("unable to marshal FlareSolverr payload to JSON: " + err.Error())
	}
	body, err := GetByAPIJson("POST", flareSolverrEndpoint, payload)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
