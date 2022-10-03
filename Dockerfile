# syntax=docker/dockerfile:1

FROM golang:1.19-alpine
WORKDIR /app/github.com/mcsymiv/web-hello-world
COPY go.* .
COPY . .
RUN ls
RUN go.mod download
# RUN go build -o /web-hello-world
EXPOSE 8080
# CMD ["/web-hello-world"]
