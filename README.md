# Crawler: Collect trade per 5 seconds data from Taiwan Stock Exchange
- `twsfex-crawler-trade-per-5s` is a crawler to collect [trade per 5 seconds data from Taiwan Stock Exchange](http://www.twse.com.tw/zh/page/trading/exchange/MI_5MINS.html).
- Using Golang and Docker.
- Using [Go Modules: twsfex-model](https://github.com/neofelisho/twsfex-model) to sync data model with back-end api.

## Installation
### Execute Golang code directly

Before starting, make sure the settings match the environment:
```go
func getEnvironments() {
	if dataSource = os.Getenv("dataSource"); dataSource == "" {
		dataSource = "http://www.twse.com.tw/en/exchangeReport/MI_5MINS?response=csv&date="
	}
	if apiUrl = os.Getenv("apiUrl"); apiUrl == "" {
		apiUrl = "http://mongo-api:8080/daily"
	}
}
```
Then we could test it:
```console
>go run .\main.go
```

### Execute by Docker

First, build the docker image. In this example, we use `test-crawler` as image name:
```console
docker build -t test-crawler .
```
Then run this image, remember attach to the same network if we run our mongodb, mongo-api on the specific docker network. For example, if we run them by docker-swarm (on multiple docker machines), or run on a testing environment by docker-compose. Check it by using command ```docker network ls```and ```docker container inspect mongodb```.

In this example, we use `test-network` as the `docker network` where the mongodb, mongo-api attach to.
```console
docker run --network test-network --rm test-crawler -date=20190214
```
In the script, `--rm` means auto remove the container when it exits. The `-date=20190214` specify the data date we want to collect.