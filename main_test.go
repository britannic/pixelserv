package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path"
	"reflect"
	"runtime"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOptsSetArgs(t *testing.T) {
	tTrue := true
	tFalse := false
	tStr := func(s string) *string {
		return &s
	}

	tests := []struct {
		// Test description.
		name string
		// Receiver fields.
		rFlagSet *flag.FlagSet
		rh       *bool
		rip      *string
		rport    *string
		rversion *bool
		// Parameters.
		fn func(int)
	}{
		{
			name:     "-h called",
			rFlagSet: &flag.FlagSet{},
			rh:       &tTrue,
			fn:       func(int) { return },
			rip:      tStr(""),
			rport:    tStr(""),
			rversion: &tFalse,
		},
		{
			name:     "-version called",
			rFlagSet: &flag.FlagSet{},
			rh:       &tFalse,
			fn:       func(int) { return },
			rip:      tStr(""),
			rport:    tStr(""),
			rversion: &tTrue,
		},
		{
			name:     "-ip and -port set",
			rFlagSet: &flag.FlagSet{},
			rh:       &tFalse,
			fn:       func(int) { return },
			rip:      tStr("192.168.168.1"),
			rport:    tStr("8080"),
			rversion: &tFalse,
		},
	}

	for _, tt := range tests {
		tt.rFlagSet.Usage = func() {
			fmt.Fprintf(os.Stderr, "Usage: %v [options]\n\n", path.Base(os.Args[0]))
			tt.rFlagSet.PrintDefaults()
		}
		o = &opts{
			FlagSet: tt.rFlagSet,
			help:    tt.rh,
			ip:      tt.rip,
			port:    tt.rport,
			version: tt.rversion,
		}

		exitCmd = func(int) { return }
		o.setArgs()
		Convey("Running main.setArgs() test", t, func() {
			So(*o.help, ShouldEqual, *tt.rh)
			So(*o.ip, ShouldEqual, *tt.rip)
			So(*o.port, ShouldEqual, *tt.rport)
			So(*o.version, ShouldEqual, *tt.rversion)
		})
	}
}

func TestGetOpts(t *testing.T) {
	Convey("Running main.getOpts() test", t, func() {
		act := new(bytes.Buffer)
		prog := path.Base(os.Args[0])

		want := `  -f string
    	Override default pixel with file source
  -h	Display help
  -ip string
    	IP address for ` + prog + ` to bind to (default "127.0.0.1")
  -path string
    	Set HTTP root path (default "/")
  -port string
    	Port number for ` + prog + ` to listen on (default "80")
  -version
    	Display version
`
		if runtime.GOOS == "linux" {
			want = `  -f="": Override default pixel with file source
    -h=false: Display help
    -ip="127.0.0.1": IP address for ` + prog + ` to bind to
    -path="/": Set HTTP root path (default "/")
    -port="80": Port number for ` + prog + ` to listen on
    -version=false: Display version
`
		}

		exitCmd = func(int) { return }
		origArgs := os.Args
		defer func() { os.Args = origArgs; return }()

		os.Args = []string{prog, "-convey-json", "-h"}

		o = getOpts()
		o.Init("pixelserv", flag.ContinueOnError)
		o.SetOutput(act)
		o.setArgs()

		So(fmt.Sprint(act), ShouldEqual, want)

		Convey("Now lets test with an invalid flag", func() {
			os.Args = []string{prog, "-z"}
			o = getOpts()
			o.Init("pixelserv", flag.ContinueOnError)
			o.SetOutput(act)
			o.setArgs()

			want += "flag provided but not defined: -z\n" + want + want
			So(fmt.Sprint(act), ShouldEqual, want)
		})
	})
}

func TestGetPix(t *testing.T) {
	Convey("Testing getPix() []byte, error", t, func() {
		exp := []byte{71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 255, 255, 255, 0, 0, 0, 33, 249, 4, 1, 0, 0, 0, 0, 44, 0, 0, 0, 0, 1, 0, 1, 0, 0, 2, 2, 68, 1, 0, 59}
		f := "./pix.bytes"

		err := ioutil.WriteFile(f, exp, 0644)
		So(err, ShouldBeNil)

		act, err := getPix(f)

		So(err, ShouldBeNil)
		So(string(act), ShouldEqual, string(exp))
	})
}

func TestLoadPixDefault(t *testing.T) {
	Convey("Testing LoadPix() http.HandleFunc", t, func() {
		req, err := http.NewRequest("GET", "/", nil)
		Convey("err should be nil and req not empty", func() {
			So(err, ShouldBeNil)
			So(req, ShouldNotBeEmpty)
			Convey("now lets check to see if loadPix() loads the correct content", func() {

				rr := httptest.NewRecorder()
				So(rr, ShouldNotBeEmpty)

				handler := http.HandlerFunc(loadPix)
				So(handler, ShouldNotBeEmpty)

				handler.ServeHTTP(rr, req)
				So(rr.Code, ShouldEqual, http.StatusOK)

				exp := []byte{71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 255, 255, 255, 0, 0, 0, 33, 249, 4, 1, 0, 0, 0, 0, 44, 0, 0, 0, 0, 1, 0, 1, 0, 0, 2, 2, 68, 1, 0, 59}

				b := reflect.DeepEqual(exp, rr.Body.Bytes())
				So(b, ShouldBeTrue)
			})
		})
	})
}

func TestLoadPixError(t *testing.T) {
	Convey("Testing LoadPix() http.HandleFunc with error assertion", t, func() {
		var act error
		origFile := *o.file
		*o.file = "rumpelstilzchen"

		defer func() { *o.file = origFile }()

		logFatalln = func(vals ...interface{}) {
			for _, v := range vals {
				act = v.(error)
			}
		}
		w := httptest.NewRecorder()
		r := &http.Request{}
		loadPix(w, r)
		So(act, ShouldNotBeNil)
	})
}

func TestLoadPixFile(t *testing.T) {
	f := "./pix.bytes"
	Convey("Testing LoadPix() http.HandleFunc", t, func() {
		o.Set("file", f)
		ioutil.WriteFile(f, []byte{71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 255, 255, 255, 0, 0, 0, 33, 249, 4, 1, 0, 0, 0, 0, 44, 0, 0, 0, 0, 1, 0, 1, 0, 0, 2, 2, 68, 1, 0, 59}, 0644)
		req, err := http.NewRequest("GET", "/", nil)
		Convey("err should be nil and req not empty", func() {
			So(err, ShouldBeNil)
			So(req, ShouldNotBeEmpty)
			Convey("now lets check to see if loadPix() loads the correct content", func() {

				rr := httptest.NewRecorder()
				So(rr, ShouldNotBeEmpty)

				handler := http.HandlerFunc(loadPix)
				So(handler, ShouldNotBeEmpty)

				handler.ServeHTTP(rr, req)
				So(rr.Code, ShouldEqual, http.StatusOK)

				exp := []byte{71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 255, 255, 255, 0, 0, 0, 33, 249, 4, 1, 0, 0, 0, 0, 44, 0, 0, 0, 0, 1, 0, 1, 0, 0, 2, 2, 68, 1, 0, 59}

				b := reflect.DeepEqual(exp, rr.Body.Bytes())
				So(b, ShouldBeTrue)
			})
		})
	})
}

func TestPixelServer(t *testing.T) {
	var (
		act           string
		origPixServer = pixelServer
	)

	defer func() { pixelServer = origPixServer; return }()

	Convey("Testing pixelServer()", t, func() {
		exp := "127.0.0.1:80"
		*o.httpPath = "/init"

		listenAndServe = func(s string, h http.Handler) error {
			act = s
			return nil
		}

		pixelServer(exp)
		So(act, ShouldEqual, exp)

		*o.httpPath = "/test"
		pixelServer = hearAndObey

		logFatalln = func(v ...interface{}) {
			act = fmt.Sprint(v)
			return
		}

		pixelServer("busted")
		So(act, ShouldEqual, "busted")
	})
}

var execCommand = exec.Command

func testPxSrv(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestPixelservMain", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	return cmd
}
