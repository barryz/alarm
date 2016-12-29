package g

import (
	"log"
	"runtime"
)

const (
	VERSION = "2.0.6@barryz"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
