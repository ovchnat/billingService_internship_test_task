# Download dependencies
FROM golang:1.16-alpine as deps
COPY go.mod go.sum /deps/
WORKDIR /deps
RUN go mod download

# Preparing test and build environment
FROM golang:1.16-alpine as build
COPY --from=deps /go/pkg /go/pkg
COPY . /billing
WORKDIR /billing

# Building with disabled dynamic linking so we
# can package executable in a thin scratch image
RUN go env -w CGO_ENABLED=0
RUN go env -w GOOS=linux
RUN go env -w GOARCH=amd64

# Run build
RUN go build -o ./cmd/billing ./cmd/main.go

# Package executable in a thin image
FROM scratch
COPY --from=build /billing/config /config
COPY --from=build /billing/.env /.env
COPY --from=build /billing/cmd/billing /cmd/billing
WORKDIR /cmd
CMD ["./billing"]

# Expose application port
EXPOSE 8080