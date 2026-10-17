// Command bubble-css writes the Bubble 4 tokens as CSS for a project's accent:
//
//	go run github.com/akrck02/bubble/cmd/bubble-css -accent "#F6BD60" > styles/bubble-tokens.css
//
// The accent is the main colour of the project's logo; Bubble adjusts it until it reads well.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/akrck02/bubble"
)

func main() {
	accentText := flag.String("accent", bubble.DefaultAccent.Hex(), "the main colour of the project's logo, as #RRGGBB")
	flag.Parse()

	accent, err := bubble.ParseColor(*accentText)
	if err != nil {
		fmt.Fprintln(os.Stderr, "bubble-css:", err)
		os.Exit(2)
	}
	if err := bubble.CSS(os.Stdout, accent); err != nil {
		fmt.Fprintln(os.Stderr, "bubble-css:", err)
		os.Exit(1)
	}
}
