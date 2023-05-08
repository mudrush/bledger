FROM golang:1.18-alpine3.17

# Add in extra linux commands not included in alpine image
RUN apk add --update --no-cache ca-certificates make git curl mercurial unzip

WORKDIR ./BLedger

# Copy contents
COPY . .

# Install requirements
RUN go mod download

# Build for linux
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o .bin/server ./cmd/server.go

RUN chmod +x .bin/server

# Expose external port
EXPOSE 8080

# Start server
CMD [".bin/server"]