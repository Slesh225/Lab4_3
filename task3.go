package main

import (
	"fmt"
)

// Функция проверки безопасности системы и определения порядка выполнения процессов
func isSafe(processes int, resources int, available []int, max [][]int, allocation [][]int) (bool, []int) {
	// Матрица, которая показывает, сколько ресурсов каждый процесс еще нуждается (need = max - allocation)
	need := make([][]int, processes)
	for i := range need {
		need[i] = make([]int, resources)
		for j := 0; j < resources; j++ {
			need[i][j] = max[i][j] - allocation[i][j]
		}
	}

	// Вектор, который показывает, сколько ресурсов каждого типа доступно в системе
	work := make([]int, resources)
	copy(work, available)

	// Вектор, который показывает, завершился ли каждый процесс
	finish := make([]bool, processes)

	// Порядок выполнения процессов
	order := make([]int, 0, processes)

	for {
		found := false
		for i := 0; i < processes; i++ {
			if !finish[i] {
				canAllocate := true
				for j := 0; j < resources; j++ {
					if need[i][j] > work[j] {
						canAllocate = false
						break
					}
				}
				if canAllocate {
					for j := 0; j < resources; j++ {
						work[j] += allocation[i][j]
					}
					finish[i] = true
					order = append(order, i)
					found = true
				}
			}
		}
		if !found {
			break
		}
	}

	// Проверка, все ли процессы завершились
	for i := 0; i < processes; i++ {
		if !finish[i] {
			return false, nil
		}
	}
	return true, order
}

func main() {
	// Пример параметров
	processes := 5 // Количество процессов
	resources := 3 // Количество типов ресурсов

	// Доступные ресурсы
	available := []int{3, 3, 2}

	//3 - 3 - 2
	//6 - 5 - 4
	//8 - 7 - 6
	//12 - 11 - 9
	//Дальше на всё ресурсов хватает

	// Максимальные запросы на ресурсы для каждого процесса
	max := [][]int{
		{7, 5, 3}, //0
		{3, 2, 2}, //1
		{9, 0, 2}, //2
		{2, 2, 2}, //3
		{4, 3, 3}, //4
	}

	// Текущее распределение ресурсов для каждого процесса
	allocation := [][]int{
		{0, 1, 0},
		{2, 0, 0},
		{3, 0, 2},
		{2, 1, 1},
		{0, 0, 2},
	}

	// Проверка, находится ли система в безопасном состоянии и определение порядка выполнения процессов
	safe, order := isSafe(processes, resources, available, max, allocation)
	if safe {
		fmt.Println("Система находится в безопасном состоянии.")
		fmt.Println("Порядок выполнения процессов:", order)
	} else {
		fmt.Println("Система не находится в безопасном состоянии.")
	}
}
