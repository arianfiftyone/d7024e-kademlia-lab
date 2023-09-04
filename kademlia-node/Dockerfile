FROM golang:1.21

RUN apt-get update 

RUN apt-get install iputils-ping -y

RUN apt-get install -y netcat-traditional

WORKDIR /kademlia-node

COPY ./src ./

CMD ["go", "run", "."]