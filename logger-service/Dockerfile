

FROM golang:1.23 as api
ENV CGO_ENABLED 0

RUN mkdir /app

COPY . /app/

WORKDIR /app

RUN go build -o loggerApp ./cmd/api

RUN chmod +x ./loggerApp


FROM alpine:3.20
RUN mkdir /app
COPY --from=api /app/loggerApp /app
WORKDIR /app
CMD ["./loggerApp"]


