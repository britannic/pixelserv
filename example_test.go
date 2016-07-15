package main

import "flag"

func Example_main() {
	testMainSetUp()
	main()
	//Output:

}

func testMainSetUp() {
	exitOnError = flag.ContinueOnError
	pixelServer = func(_ string) error {
		return nil
	}
}
