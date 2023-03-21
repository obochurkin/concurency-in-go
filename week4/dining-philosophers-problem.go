package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	times              int = 5
	eatTimes           int = 3
	eatingConcurrently int = 2 // no more than 2 philosophers to eat concurrently
)

type Chopstick struct {
	sync.Mutex
}

type IPhilosopher interface {
	Eat()
	releaseChopsticks()
}

type Philosopher struct {
	leftCS, rightCS *Chopstick
	number          int
}

func (p *Philosopher) eat(wg *sync.WaitGroup, ch chan<- Philosopher) {
	for i := 0; i < eatTimes; i++ {
		ch <- *p
		p.leftCS.Lock()
		p.rightCS.Lock()

		fmt.Printf("starting to eat %d \n", p.number)
		time.Sleep(100 * time.Microsecond) // putting food in the mouth
		fmt.Printf("finishing eating %d \n", p.number)

		p.releaseChopsticks()
		if i == eatTimes-1 {
			fmt.Printf("philosopher %d already fed exiting routine \n", p.number)
			wg.Done()
		}
	}
}

func (p *Philosopher) releaseChopsticks() {
	p.rightCS.Unlock()
	p.leftCS.Unlock()
}

func chopsticksFactory() []*Chopstick {
	chopsticks := make([]*Chopstick, times)
	for i := 0; i < times; i++ {
		chopsticks[i] = new(Chopstick)
	}

	return chopsticks
}

func philosopherFactory(chopsticks []*Chopstick) []*Philosopher {
	philosophers := make([]*Philosopher, times)
	for i := 0; i < times; i++ {
		philosophers[i] = &Philosopher{number: i + 1, leftCS: chopsticks[i], rightCS: chopsticks[(i+1)%times]}
	}

	return philosophers
}

// controlling that only 2 philosophers can eat in the same time
func host(wg *sync.WaitGroup, ch <-chan Philosopher) {
	defer wg.Done()
	count := 0
	for range ch {
		count++
		if count == eatingConcurrently {
			count = 0
			time.Sleep(1000 * time.Microsecond) // pause before the next batch starts to eat
		}
	}
}

// https://en.wikipedia.org/wiki/Dining_philosophers_problem
func main() {
	chopsticks := chopsticksFactory()
	philosophers := philosopherFactory(chopsticks)
	ch := make(chan Philosopher, eatingConcurrently)
	var wg sync.WaitGroup

	fmt.Println("This is a program that simulates the dining philosopher problem.")

	go host(&wg, ch)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go philosophers[i].eat(&wg, ch)
	}

	wg.Wait()
	defer close(ch)

	fmt.Println("All the philosophers should be fed now.")
	fmt.Println("Exit.")
}
