FROM ubuntu:22.04
RUN apt update
RUN DEBIAN_FRONTEND=noninteractive apt-get install mitmproxy golang-go python3 curl --yes
WORKDIR /app
COPY ./deploy/get-pip.py .
RUN python3 ./get-pip.py && rm ./get-pip.py
RUN pip3 install requests

COPY ./app .
RUN go build .
CMD ["./main"]
