FROM golang:1.21

RUN apt-get update 
RUN apt-get install iputils-ping -y

WORKDIR /kademlia-node

COPY ./src ./

CMD ["go", "run", "."]