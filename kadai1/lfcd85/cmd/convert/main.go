package main

import (
	"flag"

	"github.com/lfcd85/dojo5/kadai1/lfcd85/imgconv"
)

func main() {
	flag.Parse()
	dirName := flag.Arg(0)

	imgconv.Convert(dirName)
}
