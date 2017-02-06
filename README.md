# varnishlogbeat

varnishlogbeat collects log data from a Varnish Shared Memory file and ships it
to Elasticsearch.

varnishlogbeat uses [vago](https://github.com/phenomenes/vago).

## Status

varnishlogbeat is currently in beta but it's functional. If you encounter
performance issues or any bugs please create an issue or send a pull request.

### Requirements

* [Golang](https://golang.org/dl/) 1.7
* pkg-config
* [varnish-dev](http://www.varnish-cache.org/releases/) >= 4.1

You will also need to set `PKG_CONFIG_PATH` to the directory where
`varnishapi.pc` is located before running `go get`. For example:

```
export PKG_CONFIG_PATH=/usr/lib/pkgconfig
```

### Build

```
go get github.com/phenomenes/varnishlogbeat
cd $GOPATH/src/github.com/phenomenes/varnishlogbeat
go build .
```

### Run

Install and run [elasticsearch](https://github.com/elastic/elasticsearch).

Run `varnishlogbeat` with debugging output enabled:

```
./varnishlogbeat -c varnishlogbeat.yml -e -d "*"
```

Additionally you can install [kibana](https://github.com/elastic/kibana) to visualise the data.
