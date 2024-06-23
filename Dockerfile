FROM ubuntu:22.04
RUN apt update
RUN DEBIAN_FRONTEND=noninteractive apt-get install mitmproxy golang-go python3 curl --yes
WORKDIR /app
COPY ./deploy/get-pip .
RUN python3 ./get-pip && rm ./get-pip
RUN pip3 install requests

COPY ./app .
RUN go build .
CMD ["./main"]
