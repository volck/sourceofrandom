package cmd

import (
	"encoding/binary"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))

}

func getSize(filename string) int64 {
	stat, err := os.Stat(filename)
	if err != nil {
		fmt.Println("error stat:", err)
	}

	return stat.Size()
}

func isMaxSize(maxsize int64, filename string) bool {
	if maxsize != 0 {
		if getSize(filename) >= maxsize {
			fmt.Println("Reached maxsize")
			return true
		}
	}
	return false

}

func createandDefer(filename string) *os.File {
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
	}
	return f
}

func ismaxLength(maxint int64, i int) bool {
	if maxint != 0 {
		if i == int(maxint) {
			return true
		}
	}
	return false
}

func convertSize(size string) int64 {
	re := regexp.MustCompile("[0-9]+")
	out := re.FindAllString(size, -1)
	i, _ := strconv.ParseInt(out[0], 10, 64)
	finalsize := i * 1024 * 1024
	return finalsize
}

func printStats(i int, draw string, filename string) {
	if i%1000000 == 0 {
		fmt.Printf("#%d %s \n", i, draw)

	} else if int64(i)%getSize(filename) == 1024*1024 {
		fmt.Printf("have written %d mb", getSize(filename))
	}

}

func checkFileSizeAndFileLength(maxsize string, filename string, maxint int64, i int) {
	if maxsize != "" {
		finalMaxsize := convertSize(maxsize)
		if isMaxSize(finalMaxsize, filename) {
			if diehardheader {
				insertDiehardHeader(i)
			}
			os.Exit(0)
		}
	} else if ismaxLength(int64(maxint), i) {
		if diehardheader {
			insertDiehardHeader(i)
		}
		os.Exit(0)
	}
}

func insertDiehardHeader(i int) string {

	headerTemplate := fmt.Sprintf(`#==================================================================
# sourceofrandom generator
#==================================================================
type: d
count: %d
numbit: 32
`, i)
	return headerTemplate
}
