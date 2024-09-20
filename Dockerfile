FROM golang:1.22.6-alpine

WORKDIR /app

COPY src/ /app/

RUN go build -o bin .

ENV token ""
ENV guild ""
ENV rmcmd "true"

ENTRYPOINT ["/app/bin"]