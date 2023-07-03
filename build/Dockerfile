FROM node:17.6.0 as react-builder

WORKDIR /src

COPY ui ui

COPY makefile .

RUN make build-react

#############################################################

FROM golang:1.20 as go-builder

ARG app_version

ENV VERSION=$app_version

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

ENTRYPOINT ["./zym"]
CMD ["run"] 