# Logserver

Simple service for logging all incoming requests

Default listen address `:2000`

## Params

- `-a` or env variable `LISTEN_ADDR` for define listen address
- `-b` or env variable `RESPONSE_BODY` for define response body
 
## Usage

`
docker run -d -p 2000:2000 negasus/logserver

curl 127.0.0.1:2000
curl 127.0.0.1:2000/foo?bar=baz
`

output

`
----------[ 1 ]----------
2020-07-24 15:38:01.7058463 +0000 UTC m=+2.385672601
[172.17.0.1:41132] GET /

Host: 127.0.0.1:2000
Content-Length: 0
User-Agent: curl/7.64.1
Accept: */*


----------[ 2 ]----------
2020-07-24 15:39:22.7414747 +0000 UTC m=+83.525152201
[172.17.0.1:41136] GET /foo?bar=baz

Host: 127.0.0.1:2000
Content-Length: 0
User-Agent: curl/7.64.1
Accept: */*

`

## changelog

### v1.0.3

- add response body
- output application version on start

### v1.0.2

- add CORS headers

### v1.0.1

### v1.0.0

- initial version
