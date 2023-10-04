FROM golang:1.21

RUN apt-get update 

WORKDIR /kademlia-node

COPY /src ./

CMD ["go", "run", "."]