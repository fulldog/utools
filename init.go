package utools

import (
	"github.com/fulldog/utools/timex"
	"runtime"
	"time"
)

var osType = runtime.GOOS
var pathRoute = "/"

func init() {
	if osType == "windows" {
		pathRoute = "\\"
	}
	time.Local = timex.TimeZone
}
