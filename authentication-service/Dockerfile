

FROM golang:1.23 as api

ENV CGO_ENABLED 0

RUN mkdir /app

COPY . /app/

WORKDIR /app

RUN go build -o authApp ./cmd/api

RUN chmod +x ./authApp


FROM alpine:3.20
RUN mkdir /app
COPY --from=api /app/authApp /app
WORKDIR /app
CMD ["./authApp"]


