## Multistage build
FROM golang:1.20-alpine as build
ENV CGO_ENABLED=0

WORKDIR /src
COPY . .
RUN go mod download
RUN go build -o /app

## Multistage deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /
COPY --from=build /src/configs /configs
COPY --from=build /app /app

CMD ["/app", "server"]