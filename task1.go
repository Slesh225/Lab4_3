package main

import (
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

// Количество потоков
const numThreads = 7

// Задержка между попытками в SemaphoreSlim (в миллисекундах)
const semaphoreSlimDelay = 10

// Задержка в SpinWait (в миллисекундах)
const spinWaitDelay = 50

// Mutex для синхронизации доступа к общему ресурсу
var mutex sync.Mutex

// Semaphore для ограничения количества одновременно работающих потоков
var semaphore = make(chan struct{}, numThreads)

// SemaphoreSlim для ограничения количества одновременно работающих потоков
var semaphoreSlim = make(chan struct{}, numThreads/2)

// Barrier для синхронизации группы потоков
var barrier = make(chan struct{}, numThreads)

// SpinLock для синхронизации доступа к общему ресурсу
type SpinLock struct {
	state *int32
}

func NewSpinLock() *SpinLock {
	return &SpinLock{state: new(int32)}
}

func (s *SpinLock) Lock() {
	for !atomic.CompareAndSwapInt32(s.state, 0, 1) {
		// SpinWait
	}
}

func (s *SpinLock) Unlock() {
	atomic.StoreInt32(s.state, 0)
}

var spinLock = NewSpinLock()

// Monitor для синхронизации доступа к общему ресурсу
type Monitor struct {
	flag int32
	cond *sync.Cond
}

func NewMonitor() *Monitor {
	m := &Monitor{
		flag: 0,
	}
	m.cond = sync.NewCond(&sync.Mutex{})
	return m
}

func (m *Monitor) Lock() {
	for !atomic.CompareAndSwapInt32(&m.flag, 0, 1) {
		m.cond.L.Lock()
		m.cond.Wait()
		m.cond.L.Unlock()
	}
}

func (m *Monitor) Unlock() {
	atomic.StoreInt32(&m.flag, 0)
	m.cond.Signal()
}

var monitor = NewMonitor()

// Функция для генерации случайного ASCII символа
func generateRandomASCII() byte {
	return byte(rand.Intn(94) + 33) // ASCII символы от 33 до 126
}

// Функция для генерации символов с использованием Mutex
func generateWithMutex(wg *sync.WaitGroup, result *[]byte) {
	defer wg.Done()
	mutex.Lock()
	defer mutex.Unlock()
	symbol := generateRandomASCII()
	*result = append(*result, symbol)
	fmt.Printf("Mutex: %c\n", symbol)
}

// Функция для генерации символов с использованием Semaphore
func generateWithSemaphore(wg *sync.WaitGroup, result *[]byte) {
	defer wg.Done()
	semaphore <- struct{}{}
	defer func() { <-semaphore }()
	symbol := generateRandomASCII()
	*result = append(*result, symbol)
	fmt.Printf("Semaphore: %c\n", symbol)
}

// Функция для генерации символов с использованием SemaphoreSlim
func generateWithSemaphoreSlim(wg *sync.WaitGroup, result *[]byte) {
	defer wg.Done()
	for {
		select {
		case semaphoreSlim <- struct{}{}:
			defer func() { <-semaphoreSlim }()
			symbol := generateRandomASCII()
			*result = append(*result, symbol)
			fmt.Printf("SemaphoreSlim: %c\n", symbol)
			return
		default:
			time.Sleep(time.Duration(semaphoreSlimDelay) * time.Millisecond)
		}
	}
}

// Функция для генерации символов с использованием Barrier
func generateWithBarrier(wg *sync.WaitGroup, result *[]byte) {
	defer wg.Done()
	barrier <- struct{}{}
	defer func() { <-barrier }()
	symbol := generateRandomASCII()
	*result = append(*result, symbol)
	fmt.Printf("Barrier: %c\n", symbol)
}

// Функция для генерации символов с использованием SpinLock
func generateWithSpinLock(wg *sync.WaitGroup, result *[]byte) {
	defer wg.Done()
	spinLock.Lock()
	defer spinLock.Unlock()
	symbol := generateRandomASCII()
	*result = append(*result, symbol)
	fmt.Printf("SpinLock: %c\n", symbol)
}

// Функция для генерации символов с использованием Monitor
func generateWithMonitor(wg *sync.WaitGroup, result *[]byte) {
	defer wg.Done()
	monitor.Lock()
	defer monitor.Unlock()
	symbol := generateRandomASCII()
	*result = append(*result, symbol)
	fmt.Printf("Monitor: %c\n", symbol)
}

// Функция для генерации символов с использованием SpinWait
func generateWithSpinWait(wg *sync.WaitGroup, result *[]byte) {
	defer wg.Done()
	for {
		if len(*result) >= numThreads {
			return
		}
		symbol := generateRandomASCII()
		*result = append(*result, symbol)
		fmt.Printf("SpinWait: %c\n", symbol)
		time.Sleep(time.Duration(spinWaitDelay) * time.Millisecond)
	}
}

// Функция для измерения времени с использованием Stopwatch
func measureTime(f func(), name string) {
	start := time.Now()
	f()
	elapsed := time.Since(start)
	fmt.Printf("Время работы %s: %.3f мс\n", name, float64(elapsed.Nanoseconds())/1e6)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var wg sync.WaitGroup
	var mutexResult []byte
	var semaphoreResult []byte
	var semaphoreSlimResult []byte
	var barrierResult []byte
	var spinLockResult []byte
	var monitorResult []byte
	var spinWaitResult []byte

	// Функция для запуска генерации символов с использованием Mutex
	runMutex := func() {
		for i := 0; i < numThreads; i++ {
			wg.Add(1)
			go generateWithMutex(&wg, &mutexResult)
		}
		wg.Wait()
	}

	// Функция для запуска генерации символов с использованием Semaphore
	runSemaphore := func() {
		for i := 0; i < numThreads; i++ {
			wg.Add(1)
			go generateWithSemaphore(&wg, &semaphoreResult)
		}
		wg.Wait()
	}

	// Функция для запуска генерации символов с использованием SemaphoreSlim
	runSemaphoreSlim := func() {
		for i := 0; i < numThreads; i++ {
			wg.Add(1)
			go generateWithSemaphoreSlim(&wg, &semaphoreSlimResult)
		}
		wg.Wait()
	}

	// Функция для запуска генерации символов с использованием Barrier
	runBarrier := func() {
		for i := 0; i < numThreads; i++ {
			wg.Add(1)
			go generateWithBarrier(&wg, &barrierResult)
		}
		wg.Wait()
	}

	// Функция для запуска генерации символов с использованием SpinLock
	runSpinLock := func() {
		for i := 0; i < numThreads; i++ {
			wg.Add(1)
			go generateWithSpinLock(&wg, &spinLockResult)
		}
		wg.Wait()
	}

	// Функция для запуска генерации символов с использованием Monitor
	runMonitor := func() {
		for i := 0; i < numThreads; i++ {
			wg.Add(1)
			go generateWithMonitor(&wg, &monitorResult)
		}
		wg.Wait()
	}

	// Функция для запуска генерации символов с использованием SpinWait
	runSpinWait := func() {
		for i := 0; i < numThreads; i++ {
			wg.Add(1)
			go generateWithSpinWait(&wg, &spinWaitResult)
		}
		wg.Wait()
	}

	// Использование Stopwatch
	measureTime(runMutex, "Mutex")
	measureTime(runSemaphore, "Semaphore")
	measureTime(runSemaphoreSlim, "SemaphoreSlim")
	measureTime(runBarrier, "Barrier")
	measureTime(runSpinLock, "SpinLock")
	measureTime(runMonitor, "Monitor")
	measureTime(runSpinWait, "SpinWait")
}
