FROM golang:1.16.4-alpine3.13 AS build_base
# Set the Current Working Directory inside the container
WORKDIR /go/src/github.com/pallabpain/one-time-secret-sharing
COPY . .
RUN go mod download
# Build the Go app
RUN go build -o otss .
RUN chmod 777 otss

# Start fresh from a smaller image
FROM alpine:3.13
RUN apk add ca-certificates
# Copy the cloud-native-go executable from the base image to our target image
COPY --from=build_base /go/src/github.com/pallabpain/one-time-secret-sharing/otss /app/otss
# This container exposes port 8000 to the outside world
EXPOSE 9090
# Run the binary program produced by `go build`
CMD ["/app/otss"]