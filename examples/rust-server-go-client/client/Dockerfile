FROM golang:1.19-buster as builder
LABEL maintainer="Jasper van Herpt <jasper.v.herpt@gmail.com>"

# Create user for the app
RUN useradd -ms /bin/bash app-user

ENV GO111MODULE=on

WORKDIR /opt/local
COPY . .

RUN apt update && apt install --yes git make ca-certificates && update-ca-certificates

RUN go build -mod mod -v -o ./ ./

FROM debian:buster-slim AS client_service

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /opt/local/client /usr/bin/app
COPY --from=builder /opt/local/certs /certs

USER app-user

CMD ["/usr/bin/app"]
