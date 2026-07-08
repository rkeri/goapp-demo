FROM golang:1.26.5-alpine AS build

ARG VERSION=dev

WORKDIR /src

COPY app/go.mod ./
COPY app/*.go ./
RUN go build -ldflags "-X main.version=${VERSION}" -o /app .

FROM alpine:3.23.5

WORKDIR /app

COPY --from=build /app .

EXPOSE 8080

ENTRYPOINT ["./app"]
