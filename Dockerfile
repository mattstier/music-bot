FROM golang:1.25

# Install native dependencies required by hraban/opus.v2
RUN apt-get update && apt-get -y install libopus-dev libopusfile-dev pkg-config \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /go/src/app

# Copy Go modules first for caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN go build -o myapp .


# Run the binary
ENTRYPOINT ["/go/src/app/myapp"]