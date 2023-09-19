FROM golang:1.21

RUN apt-get update 

RUN apt-get install iputils-ping -y

RUN apt-get install -y netcat-traditional

RUN apt-get install -y screen

WORKDIR /kademlia-node

COPY /src ./src
COPY scripts ./

RUN cd src

CMD ["go", "run", "."]