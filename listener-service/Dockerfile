

FROM golang:1.23 as api
ENV CGO_ENABLED 0

RUN mkdir /app

COPY . /app/

WORKDIR /app

RUN go build -o listenerApp ./

RUN chmod +x ./listenerApp


FROM alpine:3.20
RUN mkdir /app
COPY --from=api /app/listenerApp /app
WORKDIR /app
CMD ["./listenerApp"]


