FROM golang:1.23 as build
WORKDIR /opt/src
COPY . .
RUN groupadd -g 1000 appuser &&\
    useradd -m -u 1000 -g appuser appuser
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o /opt/app

FROM scratch as prod
LABEL org.opencontainers.image.source="https://github.com/anton-yurchenko/golang-webserver-template"
LABEL org.opencontainers.image.version="0.1.0"
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build --chown=1000:0 /opt/app /app
ENTRYPOINT [ "/app" ]
