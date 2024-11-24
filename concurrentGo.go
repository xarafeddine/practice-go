package main

import (
	"fmt"
	"sync"
	"time"
)

// SECTION 1: Basic Goroutine Concepts
func basicGoroutineExample() {
	fmt.Println("\n=== Basic Goroutine Example ===")

	// Starting a goroutine
	go func() {
		fmt.Println("Hello from goroutine!")
	}()

	// Multiple goroutines
	for i := 0; i < 3; i++ {
		i := i // Create new variable for closure
		go func() {
			fmt.Printf("Goroutine %d executing\n", i)
		}()
	}

	// Give goroutines time to execute
	time.Sleep(time.Millisecond * 100)
}

// SECTION 2: Goroutine Synchronization with WaitGroups
func waitGroupExample() {
	fmt.Println("\n=== WaitGroup Example ===")

	var wg sync.WaitGroup

	// Launch multiple goroutines
	for i := 0; i < 3; i++ {
		wg.Add(1) // Increment counter
		i := i

		go func() {
			defer wg.Done() // Decrement counter when done
			fmt.Printf("Task %d completed\n", i)
		}()
	}

	wg.Wait() // Wait for all goroutines to complete
	fmt.Println("All tasks completed")
}

// SECTION 3: Basic Channel Operations
func basicChannelExample() {
	fmt.Println("\n=== Basic Channel Example ===")

	// Create an unbuffered channel
	ch := make(chan string)

	// Sender goroutine
	go func() {
		ch <- "Hello through channel!"
	}()

	// Receive value
	message := <-ch
	fmt.Println("Received:", message)

	// Buffered channel
	bufferedCh := make(chan int, 2)
	bufferedCh <- 1 // Won't block
	bufferedCh <- 2 // Won't block
	// bufferedCh <- 3 // Would block (buffer full)

	fmt.Println("Buffered channel values:", <-bufferedCh, <-bufferedCh)
}

// SECTION 4: Channel Direction
func channelDirectionExample() {
	fmt.Println("\n=== Channel Direction Example ===")

	ch := make(chan int)

	// Sender-only channel
	go sender(ch)

	// Receiver-only channel
	go receiver(ch)

	time.Sleep(time.Millisecond * 100)
}

func sender(ch chan<- int) { // Send-only channel
	for i := 0; i < 3; i++ {
		ch <- i
	}
}

func receiver(ch <-chan int) { // Receive-only channel
	for i := 0; i < 3; i++ {
		fmt.Println("Received:", <-ch)
	}
}

// SECTION 5: Select Statement
func selectExample() {
	fmt.Println("\n=== Select Example ===")

	ch1 := make(chan string)
	ch2 := make(chan string)
	done := make(chan bool)

	// Sender goroutines
	go func() {
		time.Sleep(time.Millisecond * 50)
		ch1 <- "Message from channel 1"
	}()

	go func() {
		time.Sleep(time.Millisecond * 100)
		ch2 <- "Message from channel 2"
	}()

	// Timeout and channel selection
	go func() {
		for {
			select {
			case msg1 := <-ch1:
				fmt.Println(msg1)
			case msg2 := <-ch2:
				fmt.Println(msg2)
			case <-time.After(time.Millisecond * 150):
				fmt.Println("Timeout!")
				done <- true
				return

			}
			fmt.Println("default")

		}
	}()

	<-done
}

// SECTION 6: Advanced Patterns
func advancedPatternsExample() {
	fmt.Println("\n=== Advanced Patterns Example ===")

	// Worker Pool Pattern
	jobs := make(chan int, 5)
	results := make(chan int, 5)

	// Start workers
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	// Send jobs
	for j := 1; j <= 5; j++ {
		jobs <- j
	}
	close(jobs)

	// Collect results
	for i := 1; i <= 5; i++ {
		<-results
	}
}

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Printf("Worker %d processing job %d\n", id, j)
		time.Sleep(time.Millisecond * 100)
		results <- j * 2
	}
}

// SECTION 7: Channel Closing and Range
func channelClosingExample() {
	fmt.Println("\n=== Channel Closing Example ===")

	ch := make(chan int)

	// Sender
	go func() {
		for i := 0; i < 3; i++ {
			ch <- i
		}
		close(ch) // Close channel when done
	}()

	// Receiver using range
	for num := range ch {
		fmt.Println("Received:", num)
	}

	// Checking if channel is closed
	_, ok := <-ch
	if !ok {
		fmt.Println("Channel is closed")
	}
}

// SECTION 8: Error Handling with Channels
func errorHandlingExample() {
	fmt.Println("\n=== Error Handling Example ===")

	results := make(chan error, 3)

	// Simulate some operations that might error
	for i := 0; i < 3; i++ {
		go func(id int) {
			if id == 1 {
				results <- fmt.Errorf("error in operation %d", id)
			} else {
				results <- nil
			}
		}(i)
	}

	// Collect and handle errors
	for i := 0; i < 3; i++ {
		if err := <-results; err != nil {
			fmt.Println("Error occurred:", err)
		} else {
			fmt.Println("Operation completed successfully")
		}
	}
}

func concurrentGo() {
	fmt.Println("Go Concurrency Crash Course")
	fmt.Println("==========================")

	// basicGoroutineExample()
	// waitGroupExample()
	// basicChannelExample()
	// channelDirectionExample()
	// selectExample()
	advancedPatternsExample()
	// channelClosingExample()
	// errorHandlingExample()
}
