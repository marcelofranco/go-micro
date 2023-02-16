FROM alpine:latest

RUN mkdir /app

COPY bin/mailerServiceApp /app
COPY templates /templates

CMD ["/app/mailerServiceApp"]