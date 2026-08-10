FROM golang:1.25.4 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o /userserver \
    ./cmd/UserServer/main.go


FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /

COPY --from=build /userserver /userserver
COPY --from=build /app/migrations /migrations
COPY --from=build /app/internal/config /internal/config

EXPOSE 50051

ENTRYPOINT ["/userserver"]
CMD ["-config=/internal/config/docker.yaml"]