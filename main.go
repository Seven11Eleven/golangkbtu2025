package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime/pprof"
	"time"
)

const filePath string = "/home/seveneleven/go/src/github.com/Seven11Eleven/golangkbtu2025/practice4_fastsearch/data/users.txt"

func FastSearch(out io.Writer) {
	file, _ := os.Open(filePath)
	defer file.Close()

	reader := bufio.NewReaderSize(file, 64*1024)
	buf := make([]byte, 64*1024)

	for {
		n, err := reader.Read(buf)
		if n > 0 {
			fmt.Fprint(out, string(buf[:n]))
		}
		if err != nil {
			break
		}
	}
}

func main() {
	cpuProfile, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(cpuProfile)
	defer pprof.StopCPUProfile()

	start := time.Now()

	for i := 0; i < 100; i++ {
		FastSearch(io.Discard)
	}

	fmt.Println("Time elapsed:", time.Since(start))
}
