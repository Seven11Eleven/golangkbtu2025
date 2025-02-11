package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"google.golang.org/protobuf/proto"
	userpb "practice4_fastsearch/data/proto"
)

const filePathPb = "./data/data.pb"

func FastSearchProto(out io.Writer) {
	file, _ := os.Open(filePathPb)
	defer file.Close()

	reader := bufio.NewReader(file)
	browsers := make(map[string]struct{}, 1000)
	fmt.Fprintln(out, "found users:")

	user := &userpb.User{}
	var email string
	var isAndroid, isMSIE bool
	var i int

	lengthBuf := make([]byte, 4)
	dataBuf := make([]byte, 1024)

	for {
		_, err := io.ReadFull(reader, lengthBuf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Ошибка чтения длины сообщения: %v", err)
			continue
		}

		length := binary.LittleEndian.Uint32(lengthBuf)

		if int(length) > len(dataBuf) {
			dataBuf = make([]byte, length)
		}

		_, err = io.ReadFull(reader, dataBuf[:length])
		if err != nil {
			log.Printf("Ошибка чтения данных Protobuf: %v", err)
			continue
		}

		user.Reset()
		err = proto.Unmarshal(dataBuf[:length], user)
		if err != nil {
			log.Printf("Ошибка анмаршалинга Protobuf: %v", err)
			continue
		}

		isAndroid = false
		isMSIE = false

		for _, browser := range user.Browsers {
			if strings.Contains(browser, "Android") {
				isAndroid = true
			} else if strings.Contains(browser, "MSIE") {
				isMSIE = true
			} else {
				continue
			}
			browsers[browser] = struct{}{}
		}

		if isAndroid && isMSIE {
			email = strings.Replace(user.Email, "@", " [at] ", -1)
			fmt.Fprintf(out, "[%d] %s <%s>\n", i, user.Name, email)
		}
		i++
	}

	fmt.Fprintf(out, "\nTotal unique browsers %d\n", len(browsers))
}
