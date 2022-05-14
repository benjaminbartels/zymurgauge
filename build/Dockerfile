# TODO: Devle does not work on ARM32
# FROM golang:1.17 as debug

# WORKDIR /go/src/github.com/benjaminbartels/zymurgauge

# COPY . .

# RUN go install github.com/go-delve/delve/cmd/dlv@latest

# EXPOSE 8080 2345

# CMD ["dlv", "debug", "github.com/benjaminbartels/zymurgauge/cmd/zym", "--listen=:2345", "--headless=true", "--api-version=2", "--log"]

#############################################################

FROM node:17.6.0 as react-builder

WORKDIR /src

COPY ui ui

COPY makefile .

RUN make build-react

#############################################################

FROM golang:1.17 as go-builder

WORKDIR /src

COPY --from=react-builder /src/ui/build ui/build

COPY cmd cmd
COPY internal internal
COPY go.mod .
COPY go.sum .
COPY makefile .
COPY ui/ui.go ui/ui.go

RUN make build-go

#############################################################

FROM scratch as production

WORKDIR /

COPY --from=go-builder /src/out/bin/zym . 
COPY --from=go-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

USER 1001:1001

EXPOSE 8080

ENTRYPOINT ["./zym"]
CMD ["run"] 