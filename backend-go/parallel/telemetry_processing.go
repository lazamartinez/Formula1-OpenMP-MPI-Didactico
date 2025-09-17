// telemetry.go - Procesamiento paralelo de telemetría
package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"math"
	"net/http"
	"runtime"
	"strconv"
	"sync"
)

func processTelemetry(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("telemetryFile")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	analysisType := r.FormValue("analysisType")

	// Leer y procesar datos CSV
	reader := csv.NewReader(file)
	var data []float64

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Convertir según el tipo de análisis
		value, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			continue
		}
		data = append(data, value)
	}

	// Procesamiento paralelo con goroutines
	numWorkers := runtime.NumCPU()
	chunkSize := len(data) / numWorkers
	results := make(chan float64, numWorkers)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		start := i * chunkSize
		end := start + chunkSize
		if i == numWorkers-1 {
			end = len(data)
		}

		go func(chunk []float64) {
			defer wg.Done()
			var result float64

			switch analysisType {
			case "max":
				result = maxValue(chunk)
			case "avg":
				result = averageValue(chunk)
			case "min":
				result = minValue(chunk)
			}

			results <- result
		}(data[start:end])
	}

	wg.Wait()
	close(results)

	// Combinar resultados
	var finalResult float64
	switch analysisType {
	case "max":
		finalResult = -1
		for res := range results {
			if res > finalResult {
				finalResult = res
			}
		}
	case "avg":
		sum, count := 0.0, 0
		for res := range results {
			sum += res
			count++
		}
		finalResult = sum / float64(count)
	case "min":
		finalResult = math.MaxFloat64
		for res := range results {
			if res < finalResult {
				finalResult = res
			}
		}
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"result":     finalResult,
		"processors": numWorkers,
		"dataPoints": len(data),
	})
}

// maxValue returns the maximum value in a slice of float64.
func maxValue(data []float64) float64 {
	if len(data) == 0 {
		return -1
	}
	max := data[0]
	for _, v := range data {
		if v > max {
			max = v
		}
	}
	return max
}

// minValue returns the minimum value in a slice of float64.
func minValue(data []float64) float64 {
	if len(data) == 0 {
		return math.MaxFloat64
	}
	min := data[0]
	for _, v := range data {
		if v < min {
			min = v
		}
	}
	return min
}

// averageValue returns the average of a slice of float64.
func averageValue(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}
