FROM ubuntu:22.04

# Set the working directory to /app
WORKDIR /app

COPY ./ /app/

# Make port 10086 available to the world outside this container
EXPOSE 10086

ENTRYPOINT ["./app"]

# To build the image, run the following command:
# docker build -t replicaset-test-image .