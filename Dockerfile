# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /trackin .

RUN apk update && apk add --no-cache tzdata
ENV TZ=America/Argentina/Buenos_Aires
ENV DEBIAN_FRONTEND=noninteractive

EXPOSE 4762

CMD [ "/trackin" ]