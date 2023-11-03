############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder
RUN apk update && apk add --no-cache vips-dev gcc musl-dev gcompat ffmpeg
WORKDIR /socialapp
COPY ./src .

# Fetch dependencies.
RUN go mod download
# Install pkger
RUN go install github.com/markbates/pkger/cmd/pkger@latest
# Build the binary.
RUN cd cmd/socialserver && go mod download && go mod tidy -v && pkger && go build -o /app/socialapp

#############################
## STEP 2 build a smaller image
#############################
FROM alpine:latest

ENV SOCIALSERVER_INSECURE_BIND_ADDRESS=0.0.0.0 \
    SOCIALSERVER_SECURE_BIND_ADDRESS=0.0.0.0 \
    SOCIALSERVER_INSECURE_BIND_PORT=8085 \
    SOCIALSERVER_SECURE_BIND_PORT=8543
    
RUN apk update && apk add --no-cache vips-dev ffmpeg
WORKDIR /app
# Copy compiled from builder.
COPY --from=builder /app/socialapp /app/socialapp
# Run the binary.
ENTRYPOINT ["/app/socialapp"]