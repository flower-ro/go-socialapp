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
    
RUN apk update && apk add --no-cache vips-dev ffmpeg
WORKDIR /app
# Copy compiled from builder.

COPY --from=builder /app/socialapp /app/socialapp
COPY --from=builder /socialapp/configs /app/configs


ENTRYPOINT ["/app/socialapp"]