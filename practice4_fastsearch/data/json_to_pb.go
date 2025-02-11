package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"log"
	"os"

	"google.golang.org/protobuf/proto"
	userpb "practice4_fastsearch/data/proto"
)

const inputPath = "./users.txt"
const outputPath = "./data.pb"

func main() {
	jsonFile, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("Ошибка открытия JSON-файла: %v", err)
	}
	defer jsonFile.Close()

	protoFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Ошибка создания файла Protobuf: %v", err)
	}
	defer protoFile.Close()

	scanner := bufio.NewScanner(jsonFile)
	for scanner.Scan() {
		var userJSON map[string]interface{}
		err := json.Unmarshal(scanner.Bytes(), &userJSON)
		if err != nil {
			log.Printf("Ошибка парсинга JSON: %v", err)
			continue
		}

		userProto := &userpb.User{
			Browsers: toStringSlice(userJSON["browsers"]),
			Email:    toString(userJSON["email"]),
			Name:     toString(userJSON["name"]),
		}

		data, _ := proto.Marshal(userProto)

		err = binary.Write(protoFile, binary.LittleEndian, uint32(len(data)))
		if err != nil {
			log.Fatalf("Ошибка записи длины сообщения: %v", err)
		}

		_, err = protoFile.Write(data)
		if err != nil {
			log.Fatalf("Ошибка записи данных Protobuf: %v", err)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Ошибка чтения JSON-файла: %v", err)
	}

	log.Println("Конвертация завершена!")
}

func toString(i interface{}) string {
	if s, ok := i.(string); ok {
		return s
	}
	return ""
}

func toStringSlice(i interface{}) []string {
	if slice, ok := i.([]interface{}); ok {
		var result []string
		for _, v := range slice {
			if s, ok := v.(string); ok {
				result = append(result, s)
			}
		}
		return result
	}
	return nil
}
