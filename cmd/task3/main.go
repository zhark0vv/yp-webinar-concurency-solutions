package main

import (
	"fmt"
	"sort"
	"sync"
)

// Job описывает задачу с приоритетом
type Job struct {
	Value    int
	Priority int
}

// worker обрабатывает задачи из канала jobs и отправляет результаты в канал results
func worker(wg *sync.WaitGroup, f func(n int) int, jobs <-chan Job, results chan<- int) {
	defer wg.Done()
	for job := range jobs {
		results <- f(job.Value)
	}
}

func main() {
	const (
		numJobs    = 10
		numWorkers = 3
	)

	// Исходные задачи
	tasks := []Job{
		{Value: 1, Priority: 4},
		{Value: 2, Priority: 3},
		{Value: 3, Priority: 2},
		{Value: 4, Priority: 1},
		{Value: 5, Priority: 0},
	}

	// Сортируем задачи по приоритету (убывание)
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Priority < tasks[j].Priority
	})

	// Создаем каналы
	jobs := make(chan Job, numJobs)
	results := make(chan int, numJobs)
	wg := sync.WaitGroup{}

	mul := func(n int) int {
		return n * 2
	}

	// Запуск воркеров
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(&wg, mul, jobs, results)
	}

	// Отправка задач в канал
	go func() {
		for _, task := range tasks {
			jobs <- task
		}
		close(jobs)
	}()

	// Закрытие канала результатов после завершения всех воркеров
	go func() {
		wg.Wait()
		close(results)
	}()

	// Вывод результатов
	fmt.Println("Results:")
	for result := range results {
		fmt.Println(result)
	}
}
