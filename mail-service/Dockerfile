

FROM golang:1.22.4-alpine3.20 as api
ENV CGO_ENABLED 0

RUN mkdir /app

COPY . /app/

WORKDIR /app

RUN go build -o mailer ./cmd/api

RUN chmod +x ./mailer


FROM alpine:3.20
RUN mkdir /app
COPY --from=api /app/mailer /app
COPY --from=api /app/template /app/template

WORKDIR /app
CMD ["./mailer"]


