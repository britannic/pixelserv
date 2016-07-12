package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"syscall"
)

var (
	// Version vars updated by go build -ldflags
	build   = "UNKNOWN"
	githash = "UNKNOWN"
	version = "UNKNOWN"

	exitCmd        = syscall.Exit
	exitOnError    = flag.ExitOnError
	listenAndServe = http.ListenAndServe
	o              *opts
)

// opts struct for command line options and setting initial variables
type opts struct {
	*flag.FlagSet
	h       *bool
	ip      *string
	port    *string
	version *bool
}

func (o *opts) setArgs(fn func(int)) {
	if os.Args[1:] != nil {
		if err := o.Parse(os.Args[1:]); err != nil {
			o.Usage()
		}
	}

	switch {
	case *o.h:
		o.Usage()
		fn(0)
	case *o.version:
		fmt.Printf(" Version:\t\t%s\n Build date:\t\t%s\n Git short hash:\t%v\n", version, build, githash)
		fn(0)
	}
}

// getOpts returns command line flags and values or displays help
func getOpts() *opts {
	flags := flag.NewFlagSet("pixelserv", exitOnError)
	flags.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %v [options]\n\n", path.Base(os.Args[0]))
		flags.PrintDefaults()
	}

	return &opts{
		h:       flags.Bool("h", false, "Display program help"),
		ip:      flags.String("ip", "127.0.0.1", "IP address for "+path.Base(os.Args[0])+" to bind to"),
		port:    flags.String("port", "80", "Port number for "+path.Base(os.Args[0])+" to listen on"),
		FlagSet: flags,
		version: flags.Bool("version", false, "Display program version number"),
	}
}

func loadPix(w http.ResponseWriter, r *http.Request) {
	pix := []byte{71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 255, 255, 255, 0, 0, 0, 33, 249, 4, 1, 0, 0, 0, 0, 44, 0, 0, 0, 0, 1, 0, 1, 0, 0, 2, 2, 68, 1, 0, 59}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "image/gif")
	w.Header().Set("Content-Length", "43")
	w.Header().Set("Accept-Ranges", "bytes")
	w.Write(pix)
}

func init() {
	o = getOpts()
}

func main() {
	o.setArgs(func(code int) {
		exitCmd(code)
	})
	http.HandleFunc("/", loadPix)
	log.Fatal(listenAndServe(fmt.Sprintf("%v:%v", *o.ip, *o.port), nil))
}
