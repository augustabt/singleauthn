# syntax=docker/dockerfile:1

FROM golang:1.19.3-alpine3.16

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd ./cmd/
COPY helpers ./helpers/
COPY html ./html/
COPY models ./models/
COPY routes ./routes/

WORKDIR /app/cmd/singleauthn

RUN go build -o /singleauthn

COPY data /data/
VOLUME "/data"

EXPOSE 7633/tcp
CMD [ "/singleauthn" ]