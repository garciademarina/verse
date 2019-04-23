# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from golang v1.12 base image
FROM golang:1.12

# Add Maintainer Info
LABEL maintainer="Carlos García de Marina Vilar <garciademarina@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/github.com/garciademarina/verse

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

WORKDIR $GOPATH/src/github.com/garciademarina/verse/cmd/api

# Download all the dependencies
# https://stackoverflow.com/questions/28031603/what-do-three-dots-mean-in-go-command-line-invocations
RUN go get -d -v ./...

# Build the package
RUN go build -v ./...
RUN cp api /opt/verse
# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["/opt/verse"]