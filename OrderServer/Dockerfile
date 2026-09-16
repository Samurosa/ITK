FROM golang:1.27.1 AS build
WORKDIR /
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /orderserver ./cmd/OrderServer/main.go


FROM alpine:3.18

RUN apk add --no-cache ca-certificates

COPY --from=build /orderserver /orderserver

COPY --from=build ./internal/config /internal/config



WORKDIR /
EXPOSE 50052
ENTRYPOINT ["/orderserver"]
CMD ["-config=/internal/config/docker.yaml"]