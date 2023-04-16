package tool

import "fmt"

// RedPrintln print msg in red
func RedPrintln(msg any) {
	fmt.Printf("\033[1;31;48m%s\033[0m\n", msg)
}

// GreenPrintln print msg in green
func GreenPrintln(msg any) {
	fmt.Printf("\033[1;32;48m%s\033[0m\n", msg)
}
