FROM golang:1.13 as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download || echo no go pks to download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=readonly -v -o server

FROM scratch

COPY --from=builder /app/server /server

CMD ["/server"]
