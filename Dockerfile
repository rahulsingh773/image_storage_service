FROM golang:latest

WORKDIR /image_service

ADD image_server ./
RUN mkdir -p ./data 

CMD ["/image_service/image_server"]