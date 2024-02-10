LOG_DIR=./logs

build:
	go build -o server.out -v ./src

build-docker:
	docker build --no-cache --network host -f ./Dockerfile . --tag main

open-last-log:
	cat $(LOG_DIR)/`ls -t $(LOG_DIR) | head -1 `

clear-logs:
	rm -rf $(LOG_DIR)/*.log
