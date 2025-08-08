FROM golang:1.24.4-bullseye

WORKDIR /app

COPY . .
RUN go mod download

RUN make build

EXPOSE 8080

CMD ["./main"]
