#!/bin/sh

varnishd \
	-a :8080 -T localhost:6082 \
	-f /etc/varnish/default.vcl \
	-s malloc,64m

./varnishlogbeat -c varnishlogbeat.yml -e -d "*"
