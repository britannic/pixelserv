package main

import (
	"bytes"
	"flag"
	"fmt"
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
		o := &opts{
			FlagSet: tt.rFlagSet,
			h:       tt.rh,
			ip:      tt.rip,
			port:    tt.rport,
			version: tt.rversion,
		}
		o.setArgs(tt.fn)
		Convey("Running main.setArgs() test", t, func() {
			So(*o.h, ShouldEqual, *tt.rh)
			So(*o.ip, ShouldEqual, *tt.rip)
			So(*o.port, ShouldEqual, *tt.rport)
			So(*o.version, ShouldEqual, *tt.rversion)
		})
	}
}

func TestGetOpts(t *testing.T) {
	want := "&{h Display program help true false}\n"
	o := getOpts()
	o.Init("pixelserv", flag.ContinueOnError)
	out := new(bytes.Buffer)
	o.SetOutput(out)
	o.Parse([]string{"-h"})
	o.setArgs(func(code int) {
		_ = code
		return
	})
	var act string
	o.Visit(func(flag *flag.Flag) {
		act += fmt.Sprintln(flag)
	})
	Convey("Running main.getOpts() test", t, func() {
		So(act, ShouldEqual, want)
	})
}

func TestLoadPix(t *testing.T) {
	Convey("Testing LoadPix() http.HandleFunc", t, func() {
		req, err := http.NewRequest("GET", "/", nil)
		Convey("err should be nil and req not empty", func() {
			So(err, ShouldBeNil)
			So(req, ShouldNotBeEmpty)
			Convey("items", func() {

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

// func TestMain(t *testing.T) {
// 	var (
// 		act     string
// 		handler interface{}
// 	)
//
// 	exitOnError = flag.ContinueOnError
// 	exitCmd = func(code int) {
// 		_ = code
// 		return
// 	}
//
// 	listenAndServe = func(s string, h http.Handler) error {
// 		act = s
// 		return nil
// 	}
//
// 	Convey("Testing main() ListenAndServe()", t, func() {
//
// 		main()
//
// 		exp := "127.0.0.1:80"
//
// 		So(act, ShouldEqual, exp)
// 		So(handler, ShouldEqual, loadPix)
// 	})
// }
