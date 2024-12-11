package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Date представляет собой структуру для хранения даты.
type Date struct {
	Day   int
	Month int
	Year  int
}

// isDateInRange проверяет, лежит ли дата date в диапазоне от start до end.
func isDateInRange(date Date, start Date, end Date) bool {
	dateToCheck := time.Date(date.Year, time.Month(date.Month), date.Day, 0, 0, 0, 0, time.UTC)
	startDate := time.Date(start.Year, time.Month(start.Month), start.Day, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(end.Year, time.Month(end.Month), end.Day, 0, 0, 0, 0, time.UTC)
	return dateToCheck.After(startDate) && dateToCheck.Before(endDate)
}

// processDatesSingleThread обрабатывает даты в однопоточном режиме.
// Возвращает количество дат, попадающих в диапазон, и сами эти даты.
func processDatesSingleThread(dates []Date, start Date, end Date) (int, []Date) {
	count := 0
	var result []Date
	for _, date := range dates {
		if isDateInRange(date, start, end) {
			count++
			result = append(result, date)
		}
	}
	return count, result
}

// processDatesMultiThread обрабатывает даты в многопоточном режиме.
// Возвращает количество дат, попадающих в диапазон, и сами эти даты.
func processDatesMultiThread(dates []Date, start Date, end Date, numThreads int) (int, []Date) {
	var count int
	var result []Date
	var wg sync.WaitGroup
	results := make(chan Date, len(dates))

	chunkSize := len(dates) / numThreads
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		startIndex := i * chunkSize
		endIndex := (i + 1) * chunkSize
		if i == numThreads-1 {
			endIndex = len(dates)
		}
		go func(chunk []Date) {
			defer wg.Done()
			for _, date := range chunk {
				if isDateInRange(date, start, end) {
					results <- date
				}
			}
		}(dates[startIndex:endIndex])
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for date := range results {
		count++
		result = append(result, date)
	}

	return count, result
}

// generateRandomDates генерирует указанное количество случайных дат в диапазоне от 01.01.1900 до 01.01.2050.
// Учитывает количество дней в каждом месяце и високосные годы.
func generateRandomDates(count int) []Date {
	dates := make([]Date, count)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < count; i++ {
		year := 1900 + rand.Intn(150)
		month := 1 + rand.Intn(12)
		day := 1 + rand.Intn(daysInMonth(year, month))
		dates[i] = Date{
			Day:   day,
			Month: month,
			Year:  year,
		}
	}
	return dates
}

// daysInMonth возвращает количество дней в месяце для указанного года.
// Учитывает високосные годы.
func daysInMonth(year, month int) int {
	switch month {
	case 2:
		if year%4 == 0 && (year%100 != 0 || year%400 == 0) {
			return 29 // Високосный год
		}
		return 28
	case 4, 6, 9, 11:
		return 30
	default:
		return 31
	}
}

func main() {
	dates100 := generateRandomDates(100)
	dates100000 := generateRandomDates(100000)

	start := Date{11, 9, 2001}
	end := Date{12, 4, 2005}
	numThreads := 4

	// Обработка базы данных на 100 дат
	startSingleThread100 := time.Now()
	countSingleThread100, resultSingleThread100 := processDatesSingleThread(dates100, start, end)
	elapsedSingleThread100 := time.Since(startSingleThread100)

	startMultiThread100 := time.Now()
	countMultiThread100, resultMultiThread100 := processDatesMultiThread(dates100, start, end, numThreads)
	elapsedMultiThread100 := time.Since(startMultiThread100)

	// Обработка базы данных на 100 000 дат
	startSingleThread100000 := time.Now()
	countSingleThread100000, _ := processDatesSingleThread(dates100000, start, end)
	elapsedSingleThread100000 := time.Since(startSingleThread100000)

	startMultiThread100000 := time.Now()
	countMultiThread100000, _ := processDatesMultiThread(dates100000, start, end, numThreads)
	elapsedMultiThread100000 := time.Since(startMultiThread100000)

	// Вывод результатов для базы данных на 100 дат
	fmt.Println("Результаты обработки базы данных на 100 дат:")
	fmt.Printf("Однопоточная обработка: %d дат, время: %.3f мс\n", countSingleThread100, float64(elapsedSingleThread100.Nanoseconds())/1e6)
	fmt.Println("Даты, подходящие под условие:")
	for _, date := range resultSingleThread100 {
		fmt.Printf("%02d.%02d.%d\n", date.Day, date.Month, date.Year)
	}

	fmt.Printf("Многопоточная обработка: %d дат, время: %.3f мс\n", countMultiThread100, float64(elapsedMultiThread100.Nanoseconds())/1e6)
	fmt.Println("Даты, подходящие под условие:")
	for _, date := range resultMultiThread100 {
		fmt.Printf("%02d.%02d.%d\n", date.Day, date.Month, date.Year)
	}

	// Вывод результатов для базы данных на 100 000 дат
	fmt.Println("Результаты обработки базы данных на 100 000 дат:")
	fmt.Printf("Однопоточная обработка: %d дат, время: %.3f мс\n", countSingleThread100000, float64(elapsedSingleThread100000.Nanoseconds())/1e6)
	fmt.Printf("Многопоточная обработка: %d дат, время: %.3f мс\n", countMultiThread100000, float64(elapsedMultiThread100000.Nanoseconds())/1e6)
}
