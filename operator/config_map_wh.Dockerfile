FROM golang:1.22-alpine as build

COPY go.mod go.sum ./
RUN go mod download
RUN go get k8s.io/client-go/tools/clientcmd@v0.30.1

COPY config_map_webhook/main.go ./main.go
COPY config_map_webhook/config_map_controller.go ./config_map_controller.go
COPY ./api ./api

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /config_map_wh .

FROM alpine:3.12

COPY --from=build /config_map_wh /config_map_wh

ENTRYPOINT ["/config_map_wh"]