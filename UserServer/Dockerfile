FROM golang:1.27.1 AS build

WORKDIR /src/UserServer

COPY UserServer/go.mod UserServer/go.sum ./
COPY exchange-common/go.mod exchange-common/go.sum /src/exchange-common/
RUN go mod download

COPY exchange-common/ /src/exchange-common/
COPY UserServer/ .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -o /userserver \
    ./cmd/UserServer/main.go


FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /

COPY --from=build /userserver /userserver
COPY --from=build /src/UserServer/migrations /migrations
COPY --from=build /src/UserServer/internal/config /internal/config

EXPOSE 50051

ENTRYPOINT ["/userserver"]
CMD ["-config=/internal/config/docker.yaml"]
