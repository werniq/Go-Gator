FROM golang:latest as build

WORKDIR .

# Copy the Go module files
COPY . .

# Download the Go module dependencies
RUN go mod download

RUN go build -o ./bin/go-gator .

EXPOSE 8080

CMD ["./bin/go-gator"]