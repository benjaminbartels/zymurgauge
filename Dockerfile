FROM golang:1.17 as debug

WORKDIR /go/src/github.com/benjaminbartels/zymurgauge

COPY . .

RUN go install github.com/go-delve/delve/cmd/dlv@latest

EXPOSE 8080 2345

CMD ["dlv", "debug", "github.com/benjaminbartels/zymurgauge/cmd/zym", "--listen=:2345", "--headless=true", "--api-version=2", "--log"]

FROM golang:1.17 as builder

WORKDIR /src

COPY . .

RUN make build

FROM scratch

WORKDIR /

COPY --from=builder /src/out/bin/zym . 
EXPOSE 8080

ENTRYPOINT ["./zym"]

