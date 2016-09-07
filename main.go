package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
)

var (
	// Version vars updated by go build -ldflags
	build   = "UNKNOWN"
	githash = "UNKNOWN"
	version = "UNKNOWN"

	exit           = os.Exit
	exitOnError    = flag.ExitOnError
	handleFunc     = http.HandleFunc
	listenAndServe = http.ListenAndServe
	logFatalln     = log.Fatalln
	o              = getOpts()
	pixelServer    = hearAndObey
	stdErr         = os.Stderr
)

// opts struct for command line options and setting initial variables
type opts struct {
	*flag.FlagSet
	file     *string
	httpPath *string
	help     *bool
	ip       *string
	port     *string
	version  *bool
}

func hearAndObey(parms string) error {
	handleFunc(*o.httpPath, loadPix)
	return listenAndServe(parms, nil)
}

func cleanArgs(args []string) []string {
	var rArgs []string
NEXT:
	for _, a := range args {
		switch {
		case strings.HasPrefix(a, "-test"):
			continue NEXT
		case strings.HasPrefix(a, "-convey"):
			continue NEXT
		default:
			rArgs = append(rArgs, a)
		}
	}
	return rArgs
}

func (o *opts) setArgs() {
	if err := o.Parse(cleanArgs((os.Args[1:]))); err != nil {
		o.Usage()
	}

	switch {
	case *o.help:
		o.Usage()
		exit(0)

	case *o.version:
		fmt.Printf(" Version:\t\t%s\n Build date:\t\t%s\n Git short hash:\t%v\n", version, build, githash)
		exit(0)
	}
}

// getOpts returns command line flags and values or displays help
func getOpts() *opts {
	flags := flag.NewFlagSet("pixelserv", exitOnError)
	flags.Usage = func() {
		fmt.Fprintf(stdErr, "Usage: %v [options]\n\n", path.Base(os.Args[0]))
		flags.PrintDefaults()
	}

	return &opts{
		file:     flags.String("f", "", "load pixel or other content from `<file>` source"),
		help:     flags.Bool("h", false, "Display help"),
		httpPath: flags.String("path", "/", "Set HTTP root path"),
		ip:       flags.String("ip", "127.0.0.1", "IP address for "+path.Base(os.Args[0])+" to bind to"),
		port:     flags.String("port", "80", "Port number for "+path.Base(os.Args[0])+" to listen on"),
		FlagSet:  flags,
		version:  flags.Bool("version", false, "Show version"),
	}
}

func getPix(f string) ([]byte, error) {
	if f != "" {
		return ioutil.ReadFile(f)
	}

	return []byte{71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 255, 255, 255, 0, 0, 0, 33, 249, 4, 1, 0, 0, 0, 0, 44, 0, 0, 0, 0, 1, 0, 1, 0, 0, 2, 2, 68, 1, 0, 59}, nil
}

func loadPix(w http.ResponseWriter, _ *http.Request) {
	pix, err := getPix(*o.file)
	if err != nil {
		logFatalln(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", http.DetectContentType(pix))
	w.Header().Set("Content-Length", fmt.Sprint(len(pix)))
	w.Header().Set("Accept-Ranges", "bytes")
	w.Write(pix)
}

func main() {
	o.setArgs()
	logFatalln(pixelServer(fmt.Sprintf("%v:%v", *o.ip, *o.port)))
}
