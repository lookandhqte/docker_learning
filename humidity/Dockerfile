FROM golang:1.18
WORKDIR /usr/src/app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -v -o /usr/local/bin/app ./...
CMD ["humidity-app"]