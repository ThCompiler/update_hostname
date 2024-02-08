package updator

import (
	"encoding/json"
	"io"
	"net/http"
)

type IP struct {
	Query string
}

const (
	ipCheckURL = "http://ip-api.com/json/"
)

func getIp() (string, error) {
	req, err := http.Get(ipCheckURL)
	if err != nil {
		return "", err
	}

	defer func() { _ = req.Body.Close() }()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return "", err
	}

	var ip IP
	if err := json.Unmarshal(body, &ip); err != nil {
		return "", err
	}

	return ip.Query, nil
}
