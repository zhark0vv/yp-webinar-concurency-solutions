package main

import (
	"fmt"
)

type Pipeline struct {
	input <-chan int
	steps []func(<-chan int) <-chan int
}

func createPipeline(input <-chan int, steps ...func(<-chan int) <-chan int) *Pipeline {
	return &Pipeline{
		input: input,
		steps: steps,
	}
}

func (p *Pipeline) Process() <-chan int {
	output := p.input
	for _, step := range p.steps {
		output = step(output) // Передаём результат предыдущего шага в следующий
	}
	return output
}

func main() {
	// Входные данные
	sales := []int{10, 20, 30, 40, 50}
	// Параметры обработки
	threshold := 30
	percent := 10
	mul := 0.3

	input := make(chan int)

	// Запуск Pipeline
	go func() {
		defer close(input)
		for _, sale := range sales {
			input <- sale
		}
	}()

	pipe := createPipeline(
		input,
		filterStep(threshold),
		increaseStep(percent),
		sumStep(),
		multiplyStep(mul),
	)

	fmt.Printf("Итог: %d\n", <-pipe.Process())

	/*
		// Шаг 1: Фильтрация данных
		Вход: 10, 20, 30, 40, 50
		Результат после фильтрации: 30, 40, 50

		// Шаг 2: Увеличение значений на процент
			30 + (30 * 10 / 100) = 33
			40 + (40 * 10 / 100) = 44
			50 + (50 * 10 / 100) = 55
		Результат после увеличения: 33, 44, 55

		// Шаг 3: Агрегация (подсчёт суммы)
		33 + 44 + 55 = 132

		// Шаг 4: Умножение на коэффициент
		132 * 0.3 = 39.6 -> 39
	*/
}

// Шаг 1: Фильтрация данных
func filterStep(threshold int) func(<-chan int) <-chan int {
	return func(input <-chan int) <-chan int {
		output := make(chan int)
		go func() {
			defer close(output)
			for value := range input {
				if value >= threshold {
					output <- value
				}
			}
		}()
		return output
	}
}

// Шаг 2: Увеличение значений на процент
func increaseStep(percent int) func(<-chan int) <-chan int {
	return func(input <-chan int) <-chan int {
		output := make(chan int)
		go func() {
			defer close(output)
			for value := range input {
				output <- value + value*percent/100
			}
		}()
		return output
	}
}

// Шаг 3: Агрегация (подсчёт суммы)
func sumStep() func(<-chan int) <-chan int {
	return func(input <-chan int) <-chan int {
		output := make(chan int)
		go func() {
			defer close(output)
			total := 0
			for value := range input {
				total += value
			}
			output <- total
		}()
		return output
	}
}

// Шаг 4: Умножение на коэффициент
func multiplyStep(mul float64) func(<-chan int) <-chan int {
	return func(input <-chan int) <-chan int {
		output := make(chan int)
		go func() {
			defer close(output)
			for value := range input {
				output <- int(float64(value) * mul)
			}
		}()
		return output
	}
}
