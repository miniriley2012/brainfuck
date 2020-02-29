package main

import (
	"flag"
	"fmt"
	"github.com/miniriley2012/brainfuck/compile"
	"io/ioutil"
	"os"
)

var (
	lang      = flag.String("c", "", "Language to compile to. Can be \"go\" or \"c\".")
	interpret = flag.Bool("i", false, "Interpret (not implemented yet).")
)

func errorln(args ...interface{}) {
	_, _ = fmt.Fprintln(os.Stderr, args...)
}

func main() {
	flag.Usage = func() {
		fmt.Printf("Usage: %s [-i] [-c language] file\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		return
	}
	*interpret = *lang == ""
	file := flag.Arg(0)
	src, err := ioutil.ReadFile(file)
	if err != nil {
		errorln("Error: " + err.Error())
		return
	}
	if *interpret {
		errorln("Interpretation is not implemented yet. Use \"-c\" instead.")
		return
	}

	var compiler compile.Compiler

	switch *lang {
	case "c":
		compiler = compile.C
	case "go":
		compiler = compile.Go
	default:
		errorln("Error: " + *lang + " is not a supported language")
		return
	}

	res, err := compiler.Compile(src)
	if err != nil {
		errorln("Error: " + err.Error())
		return
	}
	fmt.Println(string(res))
}
