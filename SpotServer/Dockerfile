FROM golang:1.27.1 AS build
WORKDIR /src/SpotServer
COPY SpotServer/go.mod SpotServer/go.sum ./
COPY exchange-common/go.mod exchange-common/go.sum /src/exchange-common/
RUN go mod download
COPY exchange-common/ /src/exchange-common/
COPY SpotServer/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o /spotserver ./cmd/SpotServer/main.go


FROM alpine:3.18

RUN apk add --no-cache ca-certificates

COPY --from=build /spotserver /spotserver

COPY --from=build /src/SpotServer/internal/config /internal/config



WORKDIR /
EXPOSE 50052
ENTRYPOINT ["/spotserver"]
CMD ["-config=/internal/config/docker.yaml"]
