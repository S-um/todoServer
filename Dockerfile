FROM golang:1.16

WORKDIR /usr/app
COPY . .

RUN go mod vendor
RUN go build -o todoServer

CMD ["./todoServer"]
