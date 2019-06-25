package main

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

// сюда писать код

func ExecutePipeline(jobs ...job) {
	l := len(jobs)
	chans := make([]chan interface{}, l)
	for i := range chans {
		chans[i] = make(chan interface{}, 100)
	}
	for i := 0; i < l-1; i++ {
		go func(i int) {
			jobs[i](chans[i], chans[i+1])
			close(chans[i+1])
		}(i)
	}
	jobs[l-1](chans[l-1], chans[l-1])
}

func SingleHash(in, out chan interface{}) {
	mu := &sync.Mutex{}
	wg := &sync.WaitGroup{}
	for v := range in {
		wg.Add(1)
		go func(v interface{}, out chan interface{}, mu *sync.Mutex, wg *sync.WaitGroup) {
			defer wg.Done()
			data := strconv.Itoa(v.(int))
			hashMd5 := make(chan string)
			go func(hashMd5 chan string, mu *sync.Mutex, data string) {
				mu.Lock()
				hashMd5 <- DataSignerMd5(data)
				mu.Unlock()
			}(hashMd5, mu, data)
			Md5 := <-hashMd5
			crc32FromHash := make(chan string)
			go func(crc32FromHash chan string, hashMd5 string) {
				crc32FromHash <- DataSignerCrc32(hashMd5)
			}(crc32FromHash, Md5)
			Crc32 := make(chan string)
			go func(Crc32 chan string, data string) {
				Crc32 <- DataSignerCrc32(data)
			}(Crc32, data)
			two := <-crc32FromHash
			one := <-Crc32
			str := one + "~" + two
			out <- str
		}(v, out, mu, wg)
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	wg1 := &sync.WaitGroup{}
	for v := range in {
		wg1.Add(1)
		go func(v interface{}, out chan interface{}, wg1 *sync.WaitGroup) {
			defer wg1.Done()
			data := v.(string)
			mu := &sync.Mutex{}
			wg := &sync.WaitGroup{}
			var group = make(map[int]string, 6)
			for th := 0; th < 6; th++ {
				wg.Add(1)
				go func(group map[int]string, th int, mu *sync.Mutex, wg *sync.WaitGroup) {
					defer wg.Done()
					str := DataSignerCrc32(strconv.Itoa(th) + data)
					mu.Lock()
					group[th] = str
					mu.Unlock()
				}(group, th, mu, wg)
			}
			wg.Wait()
			var keys []int
			for k := range group {
				keys = append(keys, k)
			}
			sort.Ints(keys)
			var str string
			for _, k := range keys {
				str += group[k]
			}
			out <- str
		}(v, out, wg1)
	}
	wg1.Wait()
}

func CombineResults(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		var res []string
		for v := range in {
			data := v.(string)
			res = append(res, data)
		}
		sort.Strings(res)
		out <- strings.Join(res, "_")
	}(wg)
	wg.Wait()
}