FROM golang:1.18 AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 go build -o /bin/bot .

FROM alpine
COPY --from=builder /bin/bot /bin/bot
ENTRYPOINT [ "/bin/bot" ]
