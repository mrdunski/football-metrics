FROM golang as test

RUN go get -u github.com/jstemmer/go-junit-report

WORKDIR /app

ADD . .

RUN go test ./... -v | go-junit-report > report.xml

FROM scratch as artifacts
COPY --from=test /app/report.xml /

FROM golang as builder

WORKDIR /football-metrics

ADD . .

RUN go build -v

FROM golang

WORKDIR /src
COPY --from=builder /football-metrics/football-metrics /src/

RUN ls /src

ENV PORT 80

ENTRYPOINT ["/src/football-metrics"]