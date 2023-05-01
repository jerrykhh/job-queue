# job-queue

### todo list

- [ ] queue data save to redis (current: ram) and load the queue data from redius
  - [ ] enqueue and dequeue use queue.Name
- [ ] UI interface
- [ ] Unit Test
- [ ] CI/CD

### .env file
```
REDIS_ADDRESS=
REDIS_PORT=
TOKEN_HASH_KEY=
ACCESS_TOKEN_DURATION=
REFRESH_TOKEN_DURATION=
GRPC_SERVER_ADDRESS=
ROOT_LOGIN_USERNAME=
ROOT_LOGIN_PWD=
```

### redis docker
```
docker run \
  -d \
  -p 6379:6379 \
  redis:alpine3.17
```