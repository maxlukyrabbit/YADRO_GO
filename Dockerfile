FROM golang:alpine AS builder

WORKDIR /build

COPY . .

RUN go build -o ./app ./src/cmd/main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/app /build/app
COPY --from=builder /build/tests/test.txt /build/tests/test.txt

CMD ["/build/app", "/build/tests/test.txt"]
