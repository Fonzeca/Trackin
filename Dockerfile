# syntax=docker/dockerfile:1

FROM golang:1.19-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /trackin .

ENV TZ=America/Argentina/Buenos_Aires

EXPOSE 4762

CMD [ "/trackin" ]