FROM alpine

RUN mkdir "/app"

COPY bin/descriptinator /app

ENTRYPOINT ["/app/descriptinator"]