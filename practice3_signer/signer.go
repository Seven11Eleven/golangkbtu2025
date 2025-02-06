package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// сюда писать код
func ExecutePipeline(jobs ...job) {
	var wg sync.WaitGroup
	in := make(chan interface{})
	for _, joba := range jobs {
		out := make(chan interface{}, 100)
		wg.Add(1)
		go func(j job, in, out chan interface{}) {
			defer wg.Done()
			defer close(out)
			j(in, out)
		}(joba, in, out)
		in = out
	}
	wg.Wait()
}

// * SingleHash считает значение crc32(data)+"~"+crc32(md5(data))
// ( конкатенация двух строк через ~), где data - то что пришло на вход
// (по сути - числа из первой функции)
// Got: 2212294583~709660146
func SingleHash(in, out chan interface{}) {
	mut := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	for data := range in {
		wg.Add(1)
		go func(data string) {
			defer wg.Done()
			crc32Hash := make(chan string)
			crc32MD5Hash := make(chan string)
			//var md5Hash string
			go func() {
				defer close(crc32Hash)
				crc32Hash <- DataSignerCrc32(data)
			}()
			go func() {
				defer close(crc32MD5Hash)

				mut.Lock()
				md5Hash := DataSignerMd5(data)
				mut.Unlock()

				crc32MD5Hash <- DataSignerCrc32(md5Hash)
			}()
			out <- fmt.Sprintf("%s~%s", <-crc32Hash, <-crc32MD5Hash)

		}(fmt.Sprintf("%v", data)) // я знаю что это затратная функция, но .(string) возвращает !ок

	}

	wg.Wait()

}

//* MultiHash считает значение crc32(th+data)) (конкатенация цифры, приведённой к строке и строки), где th=0..5 ( т.е. 6 хешей на каждое входящее значение ),
//потом берёт конкатенацию результатов в порядке расчета (0..5), где data - то что пришло на вход
// (и ушло на выход из SingleHash)

func MultiHash(in, out chan interface{}) {
	var wg sync.WaitGroup
	for data := range in {
		wg.Add(1)
		go func(data string) {
			defer wg.Done()
			results := make([]string, 6)
			var innerWg sync.WaitGroup
			for th := 0; th <= 5; th++ {
				innerWg.Add(1)
				go func(th int) {
					defer innerWg.Done()
					hash := DataSignerCrc32(fmt.Sprintf("%d%v", th, data))
					results[th] = hash
				}(th)
			}
			innerWg.Wait()
			out <- strings.Join(results, "")
		}(fmt.Sprintf("%v", data))
	}
	//out <- strings.Join(results, "")
	wg.Wait()
}

//* CombineResults получает все результаты, сортирует (https://golang.org/pkg/sort/),
//объединяет отсортированный результат через _ (символ подчеркивания) в одну строку

func CombineResults(in, out chan interface{}) {
	var res []string
	for data := range in {
		res = append(res, fmt.Sprintf("%v", data))
	}

	sort.Strings(res)
	out <- strings.Join(res, "_")

}
