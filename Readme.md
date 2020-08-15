Image Store Service: The service is used to create/delete albums and upload/delete images to/from those albums. It publishes a kafka event every time an image is uploaded or deleted.


Compile & Build binary:<br>
    ```
    go build
    ```

build docker image:<br>
    ```
    docker build -t image_server .
    ```

bring up kafka-zk:<br>
    ```
    docker-compose up -d
    ```
    
run Image Storage Service (ISS):<br>
    ```
    docker run -d -p 3000:3000 -e KAFKA_HOST=localhost -e KAFKA_TOPIC=image_server -e HOST=localhost:3000 --name image_storage_service image_server
    ```
