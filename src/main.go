package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"update_hostname/config"
)

type IP struct {
	Query string
}

func getIp() (string, error) {
	req, err := http.Get("http://ip-api.com/json/")
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

func UpdateIP(cfg *config.Config) {
	_, err := getIp()
	if err != nil {
		log.Println(err)
		return
	}

	client :=

}

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(2 * time.Hour)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case s := <-interrupt:
			log.Print(s.String())
			break
		case <-ticker.C:
			UpdateIP(cfg)
		}
	}
}
