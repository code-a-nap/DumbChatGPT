# syntax=docker/dockerfile:1

FROM golang:latest


RUN mkdir /app
RUN mkdir -p /app/public/views
RUN mkdir /app/logs
WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY public/views/*.html ./public/views/
COPY static ./static

RUN go build -o dumbChatGPT

EXPOSE 8000

CMD [ "./dumbChatGPT" ]