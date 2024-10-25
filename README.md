# Project Setup

## Clone the Repository
```bash
git clone https://github.com/auliardana/kong-jwt.git
cd kong-jwt
```

## Install Dependencies
```bash
go mod tidy
```

## Start Docker Services
```bash
cd deployments/docker
docker compose up -d
```

## Setup Kong

### 1. Setup Services

**hello-service:**

- Endpoint: `http://localhost:7777/users` [GET]
```bash
curl -i -X POST http://localhost:8001/services/ \
     --data "name=hello-service" \
     --data "url=http://hello-service:7777/users"
```

*Make sure the host in the URL matches the container name in Docker, in this case `hello-service`.*

---

**secure-service:**

- Endpoint: `http://localhost:6666/event` [POST]
```bash
curl -i -X POST http://localhost:8001/services/ \
     --data "name=secure-service-create-event" \
     --data "url=http://secure-service:6666/event"
```

- Endpoint: `http://localhost:6666/events` [GET]
```bash
curl -i -X POST http://localhost:8001/services/ \
     --data "name=secure-service-show-event" \
     --data "url=http://secure-service:6666/events"
```

---

**auth-service:**
```bash
curl -i -X POST http://localhost:8001/services/ \
     --data "name=auth-service-register" \
     --data "url=http://auth-service:8888/register"

curl -i -X POST http://localhost:8001/services/ \
     --data "name=auth-service-login" \
     --data "url=http://auth-service:8888/login"
```

### 2. Setup Routes

```bash
curl -i -X POST http://localhost:8001/services/hello-service/routes \
     --data "paths[]=/hello" \
     --data "name=hello-route"
```

*Change `service.id`, can be seen from the response when setup service.*

```bash
curl -i -X POST http://localhost:8001/services/secure-service-create-event/routes \
     --data "paths[]=/secure/event" \
     --data "name=secure-route-create"

curl -i -X POST http://localhost:8001/services/secure-service-show-event/routes \
     --data "paths[]=/secure/events" \
     --data "name=secure-route-show"

curl -i -X POST http://localhost:8001/services/auth-service-register/routes \
     --data "paths[]=/auth/register" \
     --data "name=auth-route-register"

curl -i -X POST http://localhost:8001/services/auth-service-login/routes \
     --data "paths[]=/auth/login" \
     --data "name=auth-route-login"
```

### 3. Setup JWT Authorization

```bash
curl -i -X POST http://localhost:8001/services/hello-service/plugins \
     --data "name=jwt" \
     --data "config.key_claim_name=iss" \
     --data "config.claims_to_verify[]=exp" \
     --data "config.maximum_expiration=600"
```

### 4. Setup Consumer

```bash
curl -i -X POST http://localhost:8001/consumers \
     --data "username=hello-consumer"
```

### 5. Generate JWT Credential for the Consumer

```bash
curl -i -X POST http://localhost:8001/consumers/hello-consumer/jwt \
     --data "secret=secret_key" \  # (change with secret on jwt)
     --data "key=auth-service" \  # (adjust to the issuer on jwt)
     --data "algorithm=HS256"
```

### 6. Setup Rate Limiting

```bash
curl -i -X POST http://localhost:8001/services/hello-service/plugins \
     --data "name=rate-limiting" \
     --data "config.minute=3" \
     --data "config.policy=local"
```

## Example Requests

### Register

```bash
curl -i -X POST http://localhost:8000/auth/register \
     -H "Content-Type: application/json" \
     -d '{ \
           "username": "example_username", \
           "email": "example@gmail.com", \
           "password": "password123" \
         }'
```

### Login

```bash
curl -i -X POST http://localhost:8000/auth/login \
     -H "Content-Type: application/json" \
     -d '{ \
           "email": "example@gmail.com", \
           "password": "password123" \
         }'
```

### Access Secure API

```bash
curl -i -X GET http://localhost:8000/hello \
     --header "Authorization: Bearer <your_jwt_token>"
```
