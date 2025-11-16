package main

import (
	"fmt"
	"log"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// unbuffered channellerde okayan we yazyan hokman bolmaly.
func main1() {
	ch := make(chan int)
	defer close(ch)

	go func() {
		for value := range ch {
			fmt.Println(value)
		}
	}()

	for i := 0; i < 10; i++ {
		ch <- i
	}
}

func main2() {
	ch := make(chan int)
	defer close(ch)

	// biz kanalymyzyn acykdygyny seydip barlap bileris.
	go func() {
		for {
			value, opened := <-ch
			if !opened {
				return
			}

			fmt.Println(value)
		}
	}()

	for i := 0; i < 10; i++ {
		ch <- i
	}
}

// buffered kanallarda okayan bolmasada yazyp bilyaris
func main3() {
	ch := make(chan int, 10)
	defer close(ch)

	ch <- 1
	ch <- 2

	fmt.Printf("cap = %v and length = %v", cap(ch), len(ch))
}

// bu dine okaya
func reader(ch <-chan int) {
	for value := range ch {
		fmt.Println(value)
	}
}

// bu dine yazyar
func writer(ch chan<- int) {
	for i := 0; i < 10; i++ {
		ch <- i
	}
}

// bir-birini sonsuza cenli dinlap otyrya
func main4() {
	ch := make(chan int)
	defer close(ch)

	go writer(ch)
	go reader(ch)

	_, _ = fmt.Scanln()
}

// initialize edilmedik kanaly yapmak panika doretyar.
func main5() {
	var ch chan int
	defer close(ch)
}

// yapyk kanala yazmak panika doretyar
func main6() {
	ch := make(chan int)
	close(ch)

	go writer(ch)

	_, _ = fmt.Scanln()
}

// yapyk kanaly yapmak panika doretyar
func main7() {
	ch := make(chan int)
	close(ch)

	close(ch)
}

// sul islemedi!!!
// maglumat bar bolan kanaly yapamyzda sol maglumatlary okap bolyar
func main8() {
	ch := make(chan int, 10)

	go writer(ch)

	close(ch)

	go reader(ch)

	_, _ = fmt.Scanln()
}

// yapyk kanaldan okajak bolanymyzda kanalyn default bahasyny almaly. asakdaky meselemde 0
func main9() {
	ch := make(chan int)
	close(ch)

	for value := range ch {
		fmt.Println(value)
	}
}

func main10() {
	ch := make(chan int, 10)
	defer close(ch)

	// bu 0,1,2,...,9 sanlary kanala ugratyar
	for i := 0; i < 10; i++ {
		ch <- i
	}

	// bulam 1-i kanala ugratjak bolyar yone kanal dolanlygy sebapli panika doreyar
	ch <- 1

	fmt.Printf("len=%v and cap=%v", len(ch), cap(ch))
}

// return kanallar
func Writer() <-chan int { // goroutine-1
	ch := make(chan int)
	go func() { // goroutine-2
		for i := 0; i < 10; i++ {
			ch <- i
		}
		close(ch)
	}()

	return ch
}

func Reader(ch <-chan int) { //goroutine-3
	for {
		value, opened := <-ch
		if !opened {
			return
		}

		fmt.Println(value)
	}
}

func main11() { // goroutine-4
	ch := Writer()

	go Reader(ch)

	_, _ = fmt.Scanln()
}

// select
func main() {
	first := make(chan int)
	second := make(chan int, 10)

	select {
	case first <- 1:
	case second <- 2:
	}

	var result int

	select {
	case result = <-first:
	case result = <-second:
	}

	fmt.Println(result)

	time.Sleep(time.Second)
}

// select + default
func main13() {
	ch1 := time.Tick(time.Second)
	ch2 := time.Tick(2 * time.Second)
	for {
		select {
		case <-ch1:
			fmt.Println("ch1")
		case <-ch2:
			fmt.Println("ch2")
		default:
			fmt.Println("hi")
			time.Sleep(1 * time.Second)
		}
	}
}

func main14() {
	wg := &sync.WaitGroup{}
	mx := &sync.Mutex{}
	var counter int
	mp := make(map[int]int)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			mx.Lock()
			defer mx.Unlock()
			mp[0]++
			defer wg.Done()
			counter++
		}()
	}

	counter++

	wg.Wait()

	fmt.Println(counter)
	fmt.Println(mp[0])
}

// go run -race main.go
func main15() {
	// mutex.Lock()
	// memory blocked Locked()
	// read the value
	// update the value
	// memory unblocked Unlocked()
	// mutex.Unlock()

	// but atomic
	// just update the value

	counter := 0

	for i := 0; i < 100; i++ {
		go func() {
			counter++
		}()
	}

	fmt.Println(counter)
}

// atomic or lock free
func main16() {
	var counter int64
	for i := 0; i < 100; i++ {
		go func() {
			// register call and update. thread = register + stack
			atomic.AddInt64(&counter, 1)
		}()
	}

	fmt.Println("count:", counter)
}

// GMP model
// G-goroutine
// M-machine or thread
// P-processor

// syscal bolan yagdayynda thread sol syscal gutarynca ikinji goroutine isledip bilenok

// netpoller or multiplexer bular i/o yagdayynda thread ulananok onun yerine goroutine isletya

func main17() {
	// bu dine 1 thread ulanyar
	runtime.GOMAXPROCS(1)

	wg := sync.WaitGroup{}
	wg.Add(5)

	// queue structure: fifo + lifo
	// 0, 1, 2, 3 - fifo
	// 4 - lifo
	// netije: 4 0 1 2 3
	for i := 0; i < 5; i++ {
		go func(i int) {
			defer wg.Done()
			fmt.Println(i)
		}(i)
	}

	wg.Wait()
}

// sysmon

type MutexScoreBoardManager struct {
	l          sync.RWMutex // read write mutex. read kop bolan yagdayynda suny ulansak gowy. bolmasa normal mutex
	scoreboard map[string]int
}

func NewMutexScoreBoardManager() *MutexScoreBoardManager {
	return &MutexScoreBoardManager{
		scoreboard: make(map[string]int),
	}
}

func (m *MutexScoreBoardManager) Update(name string, value int) {
	m.l.Lock()
	defer m.l.Unlock()
	m.scoreboard[name] = value
}

func (m *MutexScoreBoardManager) Read(name string) (int, bool) {
	m.l.Lock()
	defer m.l.Unlock()
	value, ok := m.scoreboard[name]
	return value, ok
}

func main18() {
	msm := NewMutexScoreBoardManager()
	teams := []string{"Lions", "Tigers", "Bears"}
	var wg sync.WaitGroup
	wg.Add(len(teams))

	for _, team := range teams {
		go func() {
			defer wg.Done()
			for i := 0; i < 10; i++ {
				score, found := msm.Read(team)
				if !found {
					score = 10
				} else {
					score += len(team)
				}
				msm.Update(team, score)
			}
		}()
	}

	wg.Wait()

	for _, team := range teams {
		score, found := msm.Read(team)
		fmt.Println(team, score, found)
	}
}

// runtime.Goshed()
func main19() {
	go func() {
		for i := 1; i <= 5; i++ {
			fmt.Println("goroutine executing:", i)
			runtime.Gosched() // indiki goroutinanyn islemegine rugsat beryar. context switching etya bir goroutine-den beyleka
		}
	}()

	for i := 1; i <= 5; i++ {
		fmt.Println("main function executing", i)
		time.Sleep(10 * time.Millisecond)
	}
}

// deadlock
func main20() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		inGoroutine := 1
		ch1 <- inGoroutine // block

		fromMain := <-ch2
		fmt.Printf("fromMain = %d", fromMain)
	}()

	inMain := 2
	ch2 <- inMain // block

	fromGoroutine := <-ch1
	fmt.Printf("fromGoroutine = %d", fromGoroutine)
}

// solution deadlock
func main21() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		inGoroutine := 1
		ch1 <- inGoroutine // block

		fromMain := <-ch2
		fmt.Printf("fromMain = %d", fromMain)
	}()

	inMain := 2

	var fromGoroutine int

	select {
	case ch2 <- inMain:
	case fromGoroutine = <-ch1:
	}

	fmt.Printf("fromGoroutine = %d", fromGoroutine)
}

// infinite loop. sysmon stop berya > 10ms ulylary
func main22() {
	go func() {
		fmt.Println("Goroutine running...")
	}()

	fmt.Println("main finished")
}

// turned of case example.
func readFromTwoChannels(ch1, ch2 chan int) []int {
	var out []int

	for count := 0; count < 2; {
		select {
		case v, ok := <-ch1:
			if !ok {
				ch1 = nil // turned of case
				count++
				fmt.Println("ch1 closed.")
				continue
			}
			out = append(out, v)

		case v, ok := <-ch2:
			if !ok {
				ch2 = nil // turned of case
				count++
				fmt.Println("ch2 closed.")
				continue
			}

			out = append(out, v)
		}
	}

	return out
}

func main23() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		for i := 10; i < 100; i += 10 {
			ch1 <- i
		}
		close(ch1)
	}()

	go func() {
		for i := 20; i >= 0; i-- {
			ch2 <- i
		}
		close(ch2)
	}()

	output := readFromTwoChannels(ch1, ch2)

	fmt.Println(output)
}

func main24() {
	ch := make(chan []int)

	go func() {
		ch <- []int{1, 3, 5}
	}()

	data := <-ch

	fmt.Printf("data = %v", data)
}

// fan in
// n sany kanaldan gelen maglumaty bir kanala ugratmak
func MergeChannels[T any](chs ...<-chan T) <-chan T {
	outputCh := make(chan T)

	wg := sync.WaitGroup{}

	wg.Add(len(chs))

	for _, ch := range chs {
		go func() {
			defer wg.Done()

			for value := range ch {
				outputCh <- value
			}
		}()
	}

	go func() {
		wg.Wait()
		close(outputCh)
	}()

	return outputCh
}

func main25() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	ch3 := make(chan int)

	go func() {
		defer func() {
			close(ch1)
			close(ch2)
			close(ch3)
		}()
		for i := 0; i < 1000000; i++ {
			ch1 <- i
			ch2 <- i + 1
			ch3 <- i + 2
		}
	}()

	for value := range MergeChannels(ch1, ch2, ch3) {
		log.Println(value)
	}
}

// fan out
// esasy kanaldan gelen maglumaty n sany kanala nobat bilen ugratmak
func SplitChannel[T any](inputCh <-chan T, n int) []<-chan T {
	outputChs := make([]chan T, n)

	for i := 0; i < n; i++ {
		outputChs[i] = make(chan T)
	}

	go func() {
		idx := 0
		for value := range inputCh {
			outputChs[idx] <- value
			idx = (idx + 1) % n
		}

		for _, ch := range outputChs {
			close(ch)
		}
	}()

	resultChs := make([]<-chan T, n)
	for i := 0; i < n; i++ {
		resultChs[i] = outputChs[i]
	}

	return resultChs
}

func main26() {
	ch := make(chan int)

	go func() {
		defer close(ch)
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()

	numberOfSplittingChannels := 2

	channels := SplitChannel(ch, numberOfSplittingChannels)

	var wg sync.WaitGroup
	wg.Add(numberOfSplittingChannels)

	go func() {
		defer wg.Done()
		for value := range channels[0] {
			log.Printf("channel 1: %d", value)
		}
	}()

	go func() {
		defer wg.Done()
		for value := range channels[1] {
			log.Printf("channel 2: %d", value)
		}
	}()

	wg.Wait()
}

// tee bu dine bir kanaldan gelen maglumaty birnace kanala ugratmak

// transfrom pattern
func TransformChannel[T any](inputCh <-chan T, action func(T) T) <-chan T {
	outputCh := make(chan T)

	go func() {
		defer close(outputCh)
		for value := range inputCh {
			outputCh <- action(value)
		}
	}()

	return outputCh
}

func main27() {
	ch := make(chan int)

	go func() {
		defer close(ch)
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()

	multiply := func(v int) int {
		return v * 2
	}

	trCh := TransformChannel(ch, multiply)

	for value := range trCh {
		log.Println(value)
	}
}

// Filter pattern
func FilterChannel[T any](inputCh <-chan T, filter func(T) bool) <-chan T {
	outputCh := make(chan T)

	go func() {
		defer close(outputCh)
		for value := range inputCh {
			if filter(value) {
				outputCh <- value
			}
		}
	}()

	return outputCh
}

func main28() {
	ch := make(chan int)

	go func() {
		defer close(ch)
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()

	isEven := func(v int) bool {
		if v%2 == 0 {
			return true
		}

		return false
	}

	rCh := FilterChannel(ch, isEven)

	for value := range rCh {
		log.Println(value)
	}
}

// pipeline pattern
func generate[T any](values ...T) <-chan T {
	outputCh := make(chan T)

	go func() {
		defer close(outputCh)
		for _, value := range values {
			outputCh <- value
		}
	}()

	return outputCh
}

func process[T any](inputCh <-chan T, action func(T) T) <-chan T {
	outputCh := make(chan T)

	go func() {
		defer close(outputCh)
		for value := range inputCh {
			outputCh <- action(value)
		}
	}()

	return outputCh
}

func main29() {
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	mul := func(v int) int {
		return v * 2
	}

	for value := range process(generate(values...), mul) {
		log.Println(value)
	}
}

// semaphore
type Semaphore struct {
	tickets chan struct{}
}

func NewSemaphore(ticketsNumber int) *Semaphore {
	return &Semaphore{make(chan struct{}, ticketsNumber)}
}

func (s *Semaphore) Acquire() {
	s.tickets <- struct{}{}
}

func (s *Semaphore) Release() {
	<-s.tickets
}

func main30() {
	wg := sync.WaitGroup{}

	wg.Add(6)

	semaphore := NewSemaphore(4)

	for i := 0; i < 6; i++ {
		semaphore.Acquire()
		go func() {
			defer func() {
				wg.Done()
				semaphore.Release()
			}()

			fmt.Println("Working...")
			time.Sleep(5 * time.Second)
			fmt.Printf("Done and Exiting job %d ...\n", i)
		}()
	}

	wg.Wait()
}

// Barrier patter. Bu testcontainerlerde ulanylyar esasy.
// Hemmesi bile bir proses yerine yetiryaler. Concurrent is bolanok. Hemmesi hokman biri birine garasmaly
type Barrier struct {
	size       int
	beforeChan chan struct{}
	afterChan  chan struct{}
	mutex      sync.Mutex
	count      int
}

func NewBarrier(size int) *Barrier {
	return &Barrier{
		size:       size,
		beforeChan: make(chan struct{}, size),
		afterChan:  make(chan struct{}, size),
	}
}

func (b *Barrier) Before() {
	b.mutex.Lock()
	b.count++

	if b.count == b.size {
		for i := 0; i < b.size; i++ {
			b.beforeChan <- struct{}{}
		}
	}

	b.mutex.Unlock()
	<-b.beforeChan
}

func (b *Barrier) After() {
	b.mutex.Lock()
	b.count--

	if b.count == 0 {
		for i := 0; i < b.size; i++ {
			b.afterChan <- struct{}{}
		}
	}

	b.mutex.Unlock()
	<-b.afterChan
}

func main31() {
	var wg sync.WaitGroup
	wg.Add(3)

	bootstrap := func() {
		log.Println("bootstrap")
	}

	work := func() {
		log.Println("done")
	}

	count := 3

	barrier := NewBarrier(count)

	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < count; j++ {
				barrier.Before()

				bootstrap()

				barrier.After()

				work()
			}
		}()
	}

	wg.Wait()
}

// promise pattern. Errorlary handlerde check etmekde ulansak bolar.
type result[T any] struct {
	val T
	err error
}

type Promise[T any] struct {
	resultCh chan result[T]
}

func NewPromise[T any](asyncFn func() (T, error)) *Promise[T] {
	promise := &Promise[T]{
		resultCh: make(chan result[T]),
	}

	go func() {
		defer close(promise.resultCh)

		val, err := asyncFn()
		promise.resultCh <- result[T]{val: val, err: err}
	}()

	return promise
}

func (p *Promise[T]) Then(successFn func(T), errorFn func(error)) {
	go func() {
		res := <-p.resultCh
		if res.err == nil {
			successFn(res.val)
		} else {
			errorFn(res.err)
		}
	}()
}

func main32() {
	asyncJob := func() (string, error) {
		time.Sleep(1 * time.Second)
		return "ok", nil
	}

	promise := NewPromise(asyncJob)

	promise.Then(
		func(value string) {
			fmt.Println("success", value)
		},
		func(err error) {
			fmt.Println("error", err)
		},
	)

	time.Sleep(2 * time.Second)
}
