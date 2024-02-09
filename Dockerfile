FROM golang:1.19

WORKDIR /app

COPY . .

RUN apt-get clean
RUN apt-get update
RUN apt-get install jq -y

RUN make build

CMD ./server.out