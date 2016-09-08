package main

import "flag"

func Example_main() {
	testMainSetUp()
	main()
}

func testMainSetUp() {
	exitOnError = flag.ContinueOnError
	pixelServer = func(_ string) error {
		return nil
	}
}
