FROM golang:1.11-stretch

LABEL maintainer phenomenes

ENV DEBIAN_FRONTEND noninteractive
ENV VER 6.0

RUN /bin/bash -c \
	'curl -s https://packagecloud.io/install/repositories/varnishcache/varnish${VER/./}/script.deb.sh | /bin/bash' \
	&& apt-get install -y \
	  libjemalloc1 \
	  pkg-config \
	  varnish \
	  varnish-dev \
	&& apt-get clean && rm -rf /var/lib/apt/lists/*

ADD . $GOPATH/src/github.com/phenomenes/varnishlogbeat

WORKDIR $GOPATH/src/github.com/phenomenes/varnishlogbeat

ADD default.vcl /etc/varnish/default.vcl
ADD docker-entrypoint.sh /docker-entrypoint.sh

RUN sed -i 's/localhost:9200/elasticsearch:9200/' \
	$GOPATH/src/github.com/phenomenes/varnishlogbeat/varnishlogbeat.yml \
	&& go build .

EXPOSE 8080

CMD /docker-entrypoint.sh
