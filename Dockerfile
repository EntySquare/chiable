From ubuntu:18.04

COPY ./chia-blockchain /chia-blockchain

WORKDIR /root/

RUN apt-get update && apt-get install python3.7-venv python3.7-distutils git lsb-release sudo python3.7-dev -y
RUN mkdir -p rplots/plots && mkdir -p rplots/cache

WORKDIR /chia-blockchain
RUN sh install.sh