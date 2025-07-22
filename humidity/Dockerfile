FROM golang:1.18

WORKDIR /app

COPY humidity/go.mod humidity/go.sum ./ 
RUN go mod download


COPY humidity/ ./
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o humidity-app
RUN chmod +x humidity-app 
RUN ls -la /app/humidity-app

CMD ["./humidity-app"]