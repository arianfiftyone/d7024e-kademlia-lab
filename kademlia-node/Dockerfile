FROM golang:1.19

RUN apt-get update 
RUN apt-get install iputils-ping

WORKDIR /kademlia-node

COPY ./src ./

CMD ["go", "run", "."]