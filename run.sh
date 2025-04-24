#!/bin/bash

cd asyncServer/

docker rm -f redis-async

docker run -d --name redis-asynq -p 6379:6379 redis

gnome-terminal -- bash -c "make run; exec bash"

cd ..

cd secure_payment_service_challenge/

make run