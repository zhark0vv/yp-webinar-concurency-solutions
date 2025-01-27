package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	ch1 := make(chan int)
	go sendValues(ch1, 1, 2, 3, 4, 5)

	ch2 := make(chan int)
	go sendValues(ch2, 6, 7, 8, 9, 10)

	ch3 := make(chan int)
	go sendValues(ch3, 11, 12, 13, 14, 15)

	//out := mergeTwo(ch1, ch2)
	//go printChan(out)

	out := mergeN(ch1, ch2, ch3)

	go printChan(out)

	time.Sleep(1 * time.Second)
}

func sendValues(ch chan int, values ...int) {
	for _, v := range values {
		ch <- v
	}
	close(ch)
}

func printChan(ch chan int) {
	for v := range ch {
		fmt.Println(v)
	}
}

func mergeTwo(ch1, ch2 chan int) chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for ch1 != nil || ch2 != nil {
			select {
			case v, ok := <-ch1:
				if ok {
					out <- v
				} else {
					ch1 = nil
				}
			case v, ok := <-ch2:
				if ok {
					out <- v
				} else {
					ch2 = nil
				}
			}
		}
	}()
	return out
}

func mergeN(chans ...chan int) chan int {
	out := make(chan int)
	var wg sync.WaitGroup
	// Функция для чтения из одного канала и записи в выходной
	worker := func(c <-chan int) {
		defer wg.Done()
		for v := range c {
			out <- v
		}
	}
	// Запускаем worker для каждого входного канала
	for _, c := range chans {
		wg.Add(1)
		go worker(c)
	}
	// Закрываем выходной канал, когда все горутины завершат работу
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}
