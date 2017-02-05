package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/phenomenes/varnishlogbeat/beater"
)

func main() {
	err := beat.Run("varnishlogbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
