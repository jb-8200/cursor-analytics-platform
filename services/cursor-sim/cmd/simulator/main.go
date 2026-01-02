// Package main provides the cursor-sim CLI application.
// cursor-sim v2 is a seed-based Cursor API simulator that generates
// synthetic usage data matching the exact Cursor Business API.
package main

import "fmt"

// Version indicates the current release of cursor-sim.
const Version = "2.0.0"

func main() {
	fmt.Printf("cursor-sim v%s\n", Version)
}
