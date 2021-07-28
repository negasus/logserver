# Logserver

Simple service for logging all incoming requests

## Params

- `-a` or env variable `LISTEN_ADDR` for define listen address (default: `:2000`)
- `-b` or env variable `RESPONSE_BODY` for define response body (default: `empty`)
- `-c` or env variable `RESPONSE_CODE` for define response status code (default: `200`)
 
## Usage

```bash
docker run -d -p 2000:2000 negasus/logserver

curl 127.0.0.1:2000

curl --header "Content-Type: application/json" \
    --request POST \
    --data '{"username":"xyz","password":"xyz"}' \
    http://127.0.0.1:2000/api/login
```

output

```
___________[ 1 ]___________
|  2021-07-28 11:31:45.939254 +0300 MSK m=+4.341261653
|  [127.0.0.1:52654] GET /
|
|  User-Agent: [curl/7.64.1]
|  Accept: [*/*]

___________[ 2 ]___________
|  2021-07-28 11:31:48.179909 +0300 MSK m=+6.581954448
|  [127.0.0.1:52655] POST /api/login
|
|  User-Agent: [curl/7.64.1]
|  Accept: [*/*]
|  Content-Type: [application/json]
|  Content-Length: [35]

{"username":"xyz","password":"xyz"}
```

## changelog

### v1.0.4

- add response status code
- decode gzip request body if needed

### v1.0.3

- add response body
- output application version on start

### v1.0.2

- add CORS headers

### v1.0.1

### v1.0.0

- initial version
