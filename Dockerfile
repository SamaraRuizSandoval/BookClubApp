# syntax=docker/dockerfile:1

FROM golang:1.24 AS build-stage

WORKDIR /app 

# Download Go modules first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Generate Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest \
 && swag init -g main.go

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /book-club-app

# Run the tests
FROM build-stage AS run-test-stage
RUN go test -v ./...

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /book-club-app /book-club-app

EXPOSE 5000
USER nonroot:nonroot
CMD ["/book-club-app"]