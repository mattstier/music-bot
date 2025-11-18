FROM golang:1.25-alpine

# Install native dependencies required by hraban/opus.v2
RUN apk add --no-cache \
    opus-dev  \
    opusfile-dev  \
    pkgconfig \
    ffmpeg \
    build-base \
    musl-dev

# installing air, the live reloading tool
RUN go install github.com/air-verse/air@latest

# Set working directory
WORKDIR /go/src/app

# Copy Go modules first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .


CMD [ "air", "-c", ".air.toml" ]

# Alternative below:

# Build the binary
# RUN go build -o myapp .

# Run the binary, alternative to CMD but overrides the live reloading
# ENTRYPOINT ["/go/src/app/myapp"]