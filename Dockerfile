FROM golang:1.8

ENV PATH=$PATH:$GOPATH/bin

RUN apt-get update && apt-get install -y \
	libjemalloc1 \
	pkg-config \
	&& apt-get clean && rm -rf /var/lib/apt/lists/*

RUN curl -O https://repo.varnish-cache.org/pkg/5.0.0/varnish_5.0.0-1_amd64.deb \
	 -O https://repo.varnish-cache.org/pkg/5.0.0/varnish-dev_5.0.0-1_amd64.deb \
        && dpkg -i varnish_5.0.0-1_amd64.deb \
	&& dpkg -i varnish-dev_5.0.0-1_amd64.deb \
	&& rm varnish_5.0.0-1_amd64.deb varnish-dev_5.0.0-1_amd64.deb

RUN go get github.com/phenomenes/varnishlogbeat \
	&& sed -i 's/localhost:9200/elasticsearch:9200/' \
	  $GOPATH/src/github.com/phenomenes/varnishlogbeat/varnishlogbeat.yml

WORKDIR $GOPATH/src/github.com/phenomenes/varnishlogbeat

COPY default.vcl /etc/varnish/default.vcl
COPY docker-entrypoint.sh /docker-entrypoint.sh

EXPOSE 8080

CMD /docker-entrypoint.sh
