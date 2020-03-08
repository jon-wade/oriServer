# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:1.13.6 as builder

# make sure unpublished dependencies are reachable
WORKDIR /app

# Fetch dependencies first; they are less susceptible to change on every build
# and will therefore be cached for speeding up the next build
COPY ./go.mod ./go.sum ./
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main

######## Start a new stage from scratch #######
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main ./

EXPOSE ${ORI_PORT}

# Run server
CMD ["./main"]
