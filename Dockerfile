# syntax=docker/dockerfile:1

FROM golang:1.24 AS build-stage

# Set destination for COPY
WORKDIR /app 

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY main.go ./
COPY docs ./docs
COPY internal ./internal
COPY migrations ./migrations

# --- Generate Swagger docs ---
RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN swag init -g main.go -o internal/docs

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /book-club-app

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /book-club-app /book-club-app

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 5000

USER nonroot:nonroot

# Run
CMD ["/book-club-app"]