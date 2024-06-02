# syntax=docker/dockerfile:1
FROM golang:1.21.10 AS build-stage
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -gcflags="all=-N -l" -o /gophermart-api ./cmd/gophermart/main.go

FROM gcr.io/distroless/base-debian11 AS build-release-stage
WORKDIR /
COPY --from=build-stage /gophermart-api /gophermart-api
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/gophermart-api"]