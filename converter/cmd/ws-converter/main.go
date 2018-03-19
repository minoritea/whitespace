package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/minoritea/whitespace/converter"
	"io"
	"log"
	"os"
)

func openSourceFile() (*os.File, error) {
	fname := flag.Arg(0)
	if fname == "" {
		log.Fatal(`Usage: ws-converter [file]`)
	}
	return os.Open(fname)
}

func main() {
	r := flag.Bool("r", false, "convert from a readable dialect")
	flag.Parse()
	var src string
	if f, err := openSourceFile(); err != nil {
		log.Fatal(err)
	} else {
		buf := bytes.NewBufferString(src)
		if _, err := io.Copy(buf, f); err != nil {
			log.Fatal(err)
		}
		src = buf.String()
	}
	if *r {
		fmt.Print(converter.FromReadable(src))
	} else {
		fmt.Print(converter.ToReadable(src))
	}
}
