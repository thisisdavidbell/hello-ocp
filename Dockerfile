FROM golang:1.10

# Set the Current Working Directory inside the container
RUN mkdir /usr/hello
WORKDIR /usr/hello

# Copy everything from the current directory to the PWD (Present Working Directory) inside the container
COPY . .

#Build
RUN go build hello-ocp.go

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the executable
CMD ["./hello-ocp"]
