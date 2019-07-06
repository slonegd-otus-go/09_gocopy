package main

import (
	"flag"
	"fmt"
)

func main() {
	from := flag.String("from", "", "copy from: path to file")
	to := flag.String("to", "", "copy to: path to file")
	offset := flag.Int("offset", 0, "offset in bytes in file copy from")
	limit := flag.Int("limit", 0, "limit in bytes for copy data, 0 means no limit")
	flag.Parse()
	fmt.Println(*from, *to, *offset, *limit)
}
