
# Weather Application

Simple API to query NWS and get today's forecast

# Build
Build a docker image
`docker build --rm -t weatherapp:v1.0 .`

# Run
Run locally
`docker run -p 8080:8080 weatherapp:v1.0`

# Query
## Forecast
Example forecast
```shell
curl --request GET \
  --url 'http://localhost:8080/v1/forecast?latitude=30.42868&longitude=-97.84273'
```

## Doc
Get OpenAPI doc
```shell
curl --request GET \
  --url http://localhost:8080/doc
```

## Health
Get Application Health
```shell
curl --request GET \
  --url http://localhost:8080/health
```

## Simulate an error response
Error response for users to test their error handling
```shell
curl --request GET \
  --url http://localhost:8080/error
```