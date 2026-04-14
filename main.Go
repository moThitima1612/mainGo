package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)


type Job struct {
	ID    int
	Score int
}

type Result struct {
	JobID int
	Grade string
	Err   error // เก็บ error ถ้ามี
}

func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		fmt.Printf("Worker %d: Processing Job %d\n", id, j.ID)
		time.Sleep(time.Millisecond * 300)

		
		if j.Score < 10 {
			results <- Result{
				JobID: j.ID,
				Err:   errors.New("invalid score: too low"),
			}
			continue // ข้ามไปทำงานถัดไป
		}

		grade := "F"
		if j.Score >= 50 {
			grade = "P"
		}

		results <- Result{JobID: j.ID, Grade: grade, Err: nil}
	}
}

func main() {
	const numJobs = 5
	const numWorkers = 2

	jobs := make(chan Job, numJobs)
	results := make(chan Result, numJobs)
	var wg sync.WaitGroup

	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go worker(w, jobs, results, &wg)
	}

	
	scores := []int{80, 5, 45, 90, 0} 
	for i, s := range scores {
		jobs <- Job{ID: i + 1, Score: s}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	
	fmt.Println("\n--- Execution Summary ---")
	for r := range results {
		if r.Err != nil {
			fmt.Printf("Job %d: ❌ Error -> %v\n", r.JobID, r.Err)
		} else {
			fmt.Printf("Job %d: ✅ Grade -> %s\n", r.JobID, r.Grade)
		}
	}
}