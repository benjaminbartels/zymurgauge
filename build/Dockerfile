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

COPY web web

COPY makefile .

RUN make build-react

#############################################################

FROM golang:1.17 as go-builder

WORKDIR /src

COPY . .

COPY --from=react-builder /src/web/build web/build

RUN make build-go

#############################################################

FROM scratch as production

WORKDIR /

COPY --from=go-builder /src/out/bin/zym . 

USER 1001:1001

EXPOSE 8080

ENTRYPOINT ["./zym"]

