# Coshkey Tree

---

## Build project

### Local

For local assembly you need to perform

```sh
$ make build # Build project
```
## Running

### For local development

```zsh
$ docker-compose up -d
```

---

## Services

### Rest:

- http://localhost:8080

```sh
[I] ➜ curl -s -X 'POST' \
  'http://localhost:8080/v1/templates' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "id": "1"
}' | jq .
{
  "code": 5,
  "message": "template not found",
  "details": []
}
```