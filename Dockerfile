FROM golang:latest
LABEL authors="kodokuus"

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/app cmd/cloud/main.go

ENV CONFIG_PATH "config/config.yaml"
ENV REFRESH_SECRET "secret"
ENV ACCESS_SECRET "secret"
ENV DATABASE_PASSWORD "pass"
CMD ["app"]