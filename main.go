package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/cheggaaa/pb/v3"

	"github.com/slonegd-otus-go/09_gocopy/internal"
)

func main() {
	from := flag.String("from", "", "путь до файла, из которого копировать")
	to := flag.String("to", "", "путь до файла, куда копировать")
	offset := flag.Int("offset", 0, "смещение в байтах от начала копируемого файла")
	limit := flag.Int("limit", 0, "максимальный размер в байтах для копирования, 0 означает без предела")
	flag.Parse()

	fromFile, err := os.Open(*from)
	if err != nil {
		fmt.Println("не могу открыть копируемый файл:", err)
		os.Exit(1)
	}
	defer fromFile.Close()

	toFile, err := os.Create(*to)
	if err != nil {
		fmt.Println("не могу создать файл, в который копировать:", err)
		os.Exit(1)
	}
	defer toFile.Close()

	if *limit == 0 {
		stat, err := fromFile.Stat()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		*limit = int(stat.Size())
	}

	bar := pb.StartNew(*limit)
	defer bar.Finish()
	bar.SetWriter(os.Stdout)

	err = internal.Process(fromFile, toFile, *offset, *limit, func(progress int) {
		bar.SetCurrent(int64(progress))
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
