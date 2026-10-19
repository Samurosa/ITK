FROM golang:1.27.1 AS build
WORKDIR /src/OrderServer
COPY OrderServer/go.mod OrderServer/go.sum ./
COPY exchange-common/go.mod exchange-common/go.sum /src/exchange-common/
RUN go mod download
COPY exchange-common/ /src/exchange-common/
COPY OrderServer/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o /orderserver ./cmd/OrderServer/main.go


FROM alpine:3.18

RUN apk add --no-cache ca-certificates

COPY --from=build /orderserver /orderserver

COPY --from=build /src/OrderServer/internal/config /internal/config



WORKDIR /
EXPOSE 50053
ENTRYPOINT ["/orderserver"]
CMD ["-config=/internal/config/docker.yaml"]
