FROM golang:1.13

WORKDIR /app
COPY . .
ENV GOOS    windows
ENV GOARCH  amd64

RUN go mod download \
 && go build -ldflags="-s -w -H windowsgui"

RUN ls -la
