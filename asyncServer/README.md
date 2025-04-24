# Secure Payments Service Async Worker

To run this you need a redis server running.

You can do it via docker using

```
docker run -d --name redis-asynq -p 6379:6379 redis
```

You can check that the server is running by pinging it with

```
docker exec -it redis redis-cli ping
```
