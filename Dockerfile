## Multistage build
FROM golang:1.20-alpine as build
ENV CGO_ENABLED=0

WORKDIR /src
COPY . .
RUN go mod download
RUN go build -o /server

## Multistage deploy
FROM gcr.io/distroless/base-debian10

WORKDIR /
COPY --from=build /src/config /config
COPY --from=build /app /app

ENTRYPOINT ["/app"]