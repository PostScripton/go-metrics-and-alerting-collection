package main

import "os"

func main() {
	os.Exit(1) // want "os.Exit is called in main package"
}
