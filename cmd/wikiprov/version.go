package main

import "fmt"

var appname string
var version string
var commit string
var date string

func getVersion() string {
	return fmt.Sprintf("%s/%s (commit: '%s' (%s))", appname, version, commit[:6], date)
}
