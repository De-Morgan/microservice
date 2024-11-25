

FROM golang:1.23 as api

ENV CGO_ENABLED 0

RUN mkdir /app

COPY . /app/

WORKDIR /app

RUN go build -o brokerApp ./cmd/api

RUN chmod +x ./brokerApp


FROM alpine:3.20
RUN mkdir /app
COPY --from=api /app/brokerApp /app
WORKDIR /app
CMD ["./brokerApp"]


