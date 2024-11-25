FROM golang:1.22.4-alpine3.20 as api
ENV CGO_ENABLED 0

RUN mkdir /app

COPY . /app/

WORKDIR /app

RUN go build -o frontApp ./cmd/web

RUN chmod +x ./frontApp


FROM alpine:3.20
RUN mkdir /app
COPY --from=api /app/frontApp /app
WORKDIR /app
CMD ["./frontApp"]

