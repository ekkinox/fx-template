FROM golang:1.20-buster as builder

WORKDIR /app
COPY .. .

RUN go mod tidy && go mod download
RUN go build -o /bin/api ./main.go

WORKDIR /bin
RUN rm -rf /app

EXPOSE 8080

CMD ["/bin/api"]