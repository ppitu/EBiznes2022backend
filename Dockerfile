FROM ubuntu:21.10

RUN apt-get update -y

RUN DEBIAN_FRONTEND=noninteractive TZ=Etc/UTC apt-get -y install tzdata

RUN apt-get install build-essential -y
RUN apt-get install wget -y
RUN apt-get install vim -y
RUN apt-get install curl -y

RUN wget https://go.dev/dl/go1.17.8.linux-amd64.tar.gz

RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.17.8.linux-amd64.tar.gz

ENV PATH=$PATH:/usr/local/go/bin

RUN echo "export PATH=$PATH" > /etc/environment

RUN useradd -ms /bin/bash ppitu
RUN adduser ppitu sudo

COPY ./ /home/ppitu/Project

USER ppitu

WORKDIR /home/ppitu/Project