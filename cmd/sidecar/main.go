package main

import (
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/skpr/fpm-metrics-adapter/internal/sidecar"
)

var (
	cliPort = kingpin.Flag("port", "Port which will respond to requests").Default(":80").String()
	cliPath = kingpin.Flag("path", "Path which responds with metrics").Default("/metrics").String()
	cliFPM  = kingpin.Flag("fpm", "Connection string for PHP FPM").Default("127.0.0.1:9000").String()
	cliFreq = kingpin.Flag("frequency", "How frequently to fresh metrics").Default("10s").Duration()
)

func main() {
	kingpin.Parse()

	err := sidecar.Start(*cliPort, *cliPath, *cliFPM, *cliFreq)
	if err != nil {
		panic(err)
	}
}
