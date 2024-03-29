FROM golang:1.17 as builder

WORKDIR /go/src

COPY ./go.mod ./go.sum ./

ENV GO111MODULE=on

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o ./app ./main.go

FROM gcr.io/distroless/static

ENV DOTENV_PATH=/go/src/.env

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/app /go/src/app

CMD ["/go/src/app"]
