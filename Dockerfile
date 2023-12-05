FROM registry.suse.com/bci/golang:1.20

RUN mkdir -p /app
# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.* ./

RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN make server

EXPOSE 8080

# Command to run the binary
CMD ["/app/webapi/webapi"]
