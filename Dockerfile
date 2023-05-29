FROM golang:1.20

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /goechoapp

EXPOSE 3000

# Run
CMD ["/goechoapp"]

LABEL name="goechoapp" \
      version="0.0.1" \
      summary="echo service" \
      description="golang echo service"
