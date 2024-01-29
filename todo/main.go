/*
Copyright © 2024 NAME HERE md3852@drexel.edu
*/
package main

import "drexel.edu/todo/cmd"

// The code in main.go was moved from here to ./cmd/root.go. This was because
// I implemented the cobra CLI and command line processor. The files are almost
// the same except with some cobra functions. I did preserve the todo: comments
// for easier navigation.
func main() {
	cmd.Execute()
}
