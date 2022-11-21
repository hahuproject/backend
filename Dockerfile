FROM golang:latest

# RUN mkdir /build

# COPY go.mod .
# COPY go.sum .
# RUN go mod download

# RUN go build -o build/main cmd/main.go

# WORKDIR /build

# CMD ["/build/main"]

# Set destination for COPY

RUN mkdir /hahu
WORKDIR /hahu

# Download Go modules
# COPY go.mod .
# COPY go.sum .

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
ADD . .
RUN go mod download

# Build
RUN go build -o app cmd/main.go

# This is for documentation purposes only.
# To actually open the port, runtime parameters
# must be supplied to the docker command.
EXPOSE 8080

# (Optional) environment variable that our dockerised
# application can make use of. The value of environment
# variables can also be set via parameters supplied
# to the docker command on the command line.
#ENV HTTP_PORT=8081

# Run
# RUN chmod +x app
CMD [ "/hahu/app" ]