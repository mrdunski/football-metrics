# football-metrics



## Test it

```shell
go test
```

## Run it

```shell
go build football-metrics
```

Then navigate to: [http endpoint](http://localhost:8080/)

## Build and run docker image

```shell
docker build -t football-metrics .
docker run -itP football-metrics
```

## Deploy the latest build to Kubernetes

```shell
helm repo add kende https://username:password@charts.dev.kende.pl/
helm install kende/football-metrics
```