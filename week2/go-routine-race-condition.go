package main

import (
    "fmt"
    "sync"
)

// global variable which will be 
var count int
var maxIterations int = 100000

func increment(wg *sync.WaitGroup) {
    defer wg.Done()
    for i := 0; i < maxIterations; i++ {
        count++
    }
}

func decrement(wg *sync.WaitGroup) {
    defer wg.Done()
    for i := 0; i < maxIterations; i++ {
        count--
    }
}

/*
The race condition in this code occurs when both increment and decrement goroutines try to access and modify the count variable at the same time.
Since they are executing concurrently, there is no guarantee about which goroutine will execute first, and therefore, whether count will be incremented or decremented first.

As a result, the final value of count becomes non-deterministic based on how the scheduler schedules the Goroutines.

To avoid this race-condition, can be used Mutex locks to ensure mutual exclusion.
*/

func main() {
    var wg sync.WaitGroup // waits goroutines to finish
    wg.Add(2) //wg must wait for 2 goroutines called wg.Done each one completed decrement this counter.

    go increment(&wg)
    go decrement(&wg)

    wg.Wait()

    fmt.Printf("Count: %d", count)
}