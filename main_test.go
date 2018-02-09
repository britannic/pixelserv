package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOptsSetArgs(t *testing.T) {
	tTrue := true
	tFalse := false
	tStr := func(s string) *string {
		return &s
	}

	// origArgs := os.Args
	// defer func() { os.Args = origArgs }()
	// os.Args = []string{path.Base(os.Args[0]), "-convey-json", "-h"}
	// prog := path.Base(os.Args[0])
	// os.Args = []string{prog, "-convey-json", "-h"}
	// cleanArgs(os.Args[1:])

	tests := []struct {
		// Test description.
		name string
		// Receiver fields.
		rFlagSet *flag.FlagSet
		rh       *bool
		rIP      *string
		rPort    *string
		rVersion *bool
		// Parameters.
		fn func(int)
	}{
		{
			name:     "-h called",
			rFlagSet: &flag.FlagSet{},
			rh:       &tTrue,
			fn:       func(int) {},
			rIP:      tStr(""),
			rPort:    tStr(""),
			rVersion: &tFalse,
		},
		{
			name:     "-version called",
			rFlagSet: &flag.FlagSet{},
			rh:       &tFalse,
			fn:       func(int) {},
			rIP:      tStr(""),
			rPort:    tStr(""),
			rVersion: &tTrue,
		},
		{
			name:     "-ip and -port set",
			rFlagSet: &flag.FlagSet{},
			rh:       &tFalse,
			fn:       func(int) {},
			rIP:      tStr("192.168.168.1"),
			rPort:    tStr("8080"),
			rVersion: &tFalse,
		},
	}

	Convey("Running main.OptsSetArgs() test", t, func() {
		for _, tt := range tests {
			Convey(fmt.Sprintf("Running main.setArgs(%s) test", tt.name), func() {
				tt.rFlagSet.Usage = func() {
					fmt.Fprintf(os.Stderr, "Usage: %v [options]\n\n", path.Base(os.Args[0]))
					tt.rFlagSet.PrintDefaults()
				}
				o = &opts{
					FlagSet: tt.rFlagSet,
					help:    tt.rh,
					ip:      tt.rIP,
					port:    tt.rPort,
					version: tt.rVersion,
				}

				exit = func(int) {}
				o.setArgs()
				So(*o.help, ShouldEqual, *tt.rh)
				So(*o.ip, ShouldEqual, *tt.rIP)
				So(*o.port, ShouldEqual, *tt.rPort)
				So(*o.version, ShouldEqual, *tt.rVersion)
			})
		}
	})
}

func TestGetOpts(t *testing.T) {
	Convey("Running main.getOpts() test", t, func() {
		act := new(bytes.Buffer)
		prog := path.Base(os.Args[0])

		want := `  -f <file>
    	load pixel or other content from <file> source
  -h	Display help
  -ip string
    	IP address for ` + prog + ` to bind to (default "127.0.0.1")
  -path string
    	Set HTTP root path (default "/")
  -port string
    	Port number for ` + prog + ` to listen on (default "80")
  -version
    	Show version
`

		if os.Getenv("DRONE") == "true" {
			want = "  -f=\"\": load pixel or other content from `<file>` source\n  -h=false: Display help\n  -ip=\"127.0.0.1\": IP address for " + prog + " to bind to\n  -path=\"/\": Set HTTP root path\n  -port=\"80\": Port number for " + prog + " to listen on\n  -version=false: Show version\n"
		}

		exit = func(int) {}
		origArgs := os.Args
		defer func() { os.Args = origArgs }()

		os.Args = []string{prog, "-convey-json", "-h"}

		o = getOpts()
		o.Init(prog, flag.ContinueOnError)
		o.SetOutput(act)
		o.setArgs()

		So(act.String(), ShouldEqual, want)

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
		var act string
		origFile := *o.file
		*o.file = "rumpelstilzchen"

		defer func() { *o.file = origFile }()

		logFatalln = func(vals ...interface{}) {
			for _, v := range vals {
				act += v.(string)
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

	defer func() { pixelServer = origPixServer }()

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
		}

		pixelServer("busted")
		So(act, ShouldEqual, "busted")
	})
}

// var execCommand = exec.Command

// func testPxSrv(command string, args ...string) *exec.Cmd {
// 	cs := []string{"-test.run=TestPixelservMain", "--", command}
// 	cs = append(cs, args...)
// 	cmd := exec.Command(os.Args[0], cs...)
// 	return cmd
// }
