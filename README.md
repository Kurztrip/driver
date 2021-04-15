# driver
Driver microservice in Go(Golang)

-To build image

docker build -t sisarmientob/driver_app

-To run docker

sudo docker run -p 80:8080 -e HOST=34.123.78.15 sisarmientob/driver_app

34.123.78.15 is the ip of the databse server