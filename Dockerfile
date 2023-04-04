FROM ubuntu:22.04
RUN apt update
RUN DEBIAN_FRONTEND=noninteractive apt-get install mitmproxy golang-go curl --yes

WORKDIR /app
COPY ./app .
RUN go build .
CMD ["./main"]
