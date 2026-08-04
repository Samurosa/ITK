FROM golang:1.25.4 AS build
WORKDIR /
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /userserver ./cmd/UserServer/main.go


FROM alpine:3.18

RUN apk add --no-cache ca-certificates

COPY --from=build /userserver /userserver

COPY --from=build ./migrations /migrations

COPY --from=build ./internal/config /internal/config

COPY --from=build .env .env



WORKDIR /
EXPOSE 50051
ENTRYPOINT ["/userserver"]