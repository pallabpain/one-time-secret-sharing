# One-Time Secret Sharing

The One-Time Secret Sharing (OTSS) is a set of simple APIs written in Go that facilitate one-time sharing of messages and secrets.

## Starting the application
You will need `docker` and `docker-compose` installed on your system.
```
git clone https://github.com/pallabpain/one-time-secret-sharing.git
cd one-time-secret-sharing
docker-compose up
```
## APIs
Checking the state of the service
```
curl -X GET http://localhost:9090/ready
```
Creating a new secret
```
curl -H "Content-Type: application/json" -X POST http://localhost:9090/secrets -d '{"message": "this is a test message", "password": "p@ssw0rd"}'

{"uuid": 924645c1-6a3d-4ab3-89d9-1c3c8aa59b49}
```
Reading the password
```
curl -H "Content-Type: application/json" -X POST http://localhost:9090/secrets/924645c1-6a3d-4ab3-89d9-1c3c8aa59b49 -d '{"password": "p@ssw0rd"}'

{"message": "this is a test message"}
```

## Go Libraries Used
- [gomodule/redigo](https://github.com/gomodule/redigo)
- [gorilla/mux](https://github.com/gorilla/mux)
- [google/uuid](https://github.com/google/uuid)