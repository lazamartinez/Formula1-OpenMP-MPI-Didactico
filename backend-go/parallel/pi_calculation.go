// pi_calculation.go - Cálculo paralelo de π
package main

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func calculatePi(w http.ResponseWriter, r *http.Request) {
	iterations, _ := strconv.Atoi(r.FormValue("iterations"))
	numThreads, _ := strconv.Atoi(r.FormValue("threads"))

	chunkSize := iterations / numThreads
	results := make(chan float64, numThreads)
	var wg sync.WaitGroup

	startTime := time.Now()

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		start := i * chunkSize
		end := start + chunkSize
		if i == numThreads-1 {
			end = iterations
		}

		go func(start, end int) {
			defer wg.Done()
			localSum := 0.0

			for j := start; j < end; j++ {
				x := (float64(j) + 0.5) / float64(iterations)
				localSum += 4.0 / (1.0 + x*x)
			}

			results <- localSum
		}(start, end)
	}

	wg.Wait()
	close(results)

	totalSum := 0.0
	for res := range results {
		totalSum += res
	}

	pi := totalSum / float64(iterations)
	duration := time.Since(startTime)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"pi":         pi,
		"time":       duration.Milliseconds(),
		"iterations": iterations,
		"threads":    numThreads,
		"error":      math.Abs(pi-math.Pi) / math.Pi * 100,
	})
}
