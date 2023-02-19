FROM golang:alpine as builder

RUN mkdir /app

ADD . /app

WORKDIR /app

# COPY db.sql /docker-entrypoint-initdb.d/

RUN go build -o main ./cmd/my-rest-api/main.go

EXPOSE 8081
CMD [ "/app/main" ]

