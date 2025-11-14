package integration

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/urfave/cli/v2"

	"github.com/denkhaus/open-notebook-cli/pkg/di"
)

// TestPerformanceStressTests tests performance under various loads
func TestPerformanceStressTests(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance stress tests")
	}

	t.Run("Response time measurement", func(t *testing.T) {
		testResponseTimes(t)
	})

	t.Run("Concurrent load testing", func(t *testing.T) {
		testConcurrentLoad(t)
	})

	t.Run("Memory usage monitoring", func(t *testing.T) {
		testMemoryUsage(t)
	})
}

// testResponseTimes measures response times for various operations
func testResponseTimes(t *testing.T) {
	// Create a fast server for response time testing
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate different response times based on endpoint
		switch {
		case strings.Contains(r.URL.Path, "/slow"):
			time.Sleep(100 * time.Millisecond)
		case strings.Contains(r.URL.Path, "/medium"):
			time.Sleep(50 * time.Millisecond)
		default:
			time.Sleep(10 * time.Millisecond) // Fast response
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"endpoint":  r.URL.Path,
			"timestamp": time.Now().Unix(),
		})
	}))
	defer server.Close()

	ctx, err := createPerformanceCLIContextWithServer(server.URL)
	require.NoError(t, err)

	injector := di.Bootstrap(ctx)
	httpClient := di.GetHTTPClient(injector)

	endpoints := []struct {
		name string
		path string
	}{
		{"Fast endpoint", "/fast"},
		{"Medium endpoint", "/medium"},
		{"Slow endpoint", "/slow"},
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint.name, func(t *testing.T) {
			const iterations = 5
			var durations []time.Duration

			for i := 0; i < iterations; i++ {
				start := time.Now()
				resp, err := httpClient.Get(context.Background(), endpoint.path)
				duration := time.Since(start)

				require.NoError(t, err)
				assert.Equal(t, 200, resp.StatusCode)

				durations = append(durations, duration)
			}

			// Calculate statistics
			var totalDuration time.Duration
			minDuration := durations[0]
			maxDuration := durations[0]

			for _, d := range durations {
				totalDuration += d
				if d < minDuration {
					minDuration = d
				}
				if d > maxDuration {
					maxDuration = d
				}
			}

			avgDuration := totalDuration / time.Duration(len(durations))

			t.Logf("✅ %s response times: avg=%v, min=%v, max=%v (%d requests)",
				endpoint.name, avgDuration, minDuration, maxDuration, len(durations))

			// Performance assertions
			if strings.Contains(endpoint.path, "fast") {
				assert.Less(t, avgDuration, 50*time.Millisecond,
					"Fast endpoint should respond quickly")
			} else if strings.Contains(endpoint.path, "slow") {
				assert.Less(t, avgDuration, 200*time.Millisecond,
					"Slow endpoint should still respond within 200ms")
			}
		})
	}
}

// testConcurrentLoad tests performance under concurrent load
func testConcurrentLoad(t *testing.T) {
	// Create server that handles concurrent requests
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Small delay to simulate processing
		time.Sleep(20 * time.Millisecond)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"request_id":  fmt.Sprintf("req-%d", time.Now().UnixNano()),
			"server_time": time.Now().Format(time.RFC3339),
		})
	}))
	defer server.Close()

	ctx, err := createPerformanceCLIContextWithServer(server.URL)
	require.NoError(t, err)

	injector := di.Bootstrap(ctx)
	httpClient := di.GetHTTPClient(injector)

	t.Run("Moderate concurrent load", func(t *testing.T) {
		const numWorkers = 5
		const requestsPerWorker = 10

		var wg sync.WaitGroup
		results := make(chan TestResult, numWorkers*requestsPerWorker)

		start := time.Now()

		// Launch workers
		for i := 0; i < numWorkers; i++ {
			wg.Add(1)
			go func(workerID int) {
				defer wg.Done()

				for j := 0; j < requestsPerWorker; j++ {
					reqStart := time.Now()
					resp, err := httpClient.Get(context.Background(), "/test")
					duration := time.Since(reqStart)

					result := TestResult{
						WorkerID: workerID,
						Success:  err == nil,
						Duration: duration,
					}

					if err != nil {
						result.Error = err.Error()
					} else if resp.StatusCode != 200 {
						result.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
					}

					results <- result
				}
			}(i)
		}

		wg.Wait()
		close(results)

		totalDuration := time.Since(start)

		// Collect and analyze results
		var successCount, errorCount int
		var totalResponseTime time.Duration
		maxDuration := time.Duration(0)
		minDuration := time.Hour

		for result := range results {
			if result.Success {
				successCount++
				totalResponseTime += result.Duration
				if result.Duration > maxDuration {
					maxDuration = result.Duration
				}
				if result.Duration < minDuration {
					minDuration = result.Duration
				}
			} else {
				errorCount++
			}
		}

		totalRequests := successCount + errorCount
		successRate := float64(successCount) / float64(totalRequests) * 100
		avgDuration := totalResponseTime / time.Duration(successCount)

		t.Logf("✅ Concurrent load test completed:")
		t.Logf("   Total requests: %d", totalRequests)
		t.Logf("   Success rate: %.1f%%", successRate)
		t.Logf("   Average response time: %v", avgDuration)
		t.Logf("   Min/Max response time: %v/%v", minDuration, maxDuration)
		t.Logf("   Total duration: %v", totalDuration)
		t.Logf("   Requests per second: %.1f", float64(totalRequests)/totalDuration.Seconds())

		// Assertions
		assert.GreaterOrEqual(t, successRate, 90.0,
			"Success rate should be at least 90%")
		assert.Less(t, avgDuration, 500*time.Millisecond,
			"Average response time should be reasonable")
	})

	t.Run("High concurrent load", func(t *testing.T) {
		const burstSize = 20
		const maxDuration = 30 * time.Second

		ctx, cancel := context.WithTimeout(context.Background(), maxDuration)
		defer cancel()

		var wg sync.WaitGroup
		results := make(chan TestResult, burstSize)

		start := time.Now()

		// Launch burst of concurrent requests
		for i := 0; i < burstSize; i++ {
			wg.Add(1)
			go func(requestID int) {
				defer wg.Done()

				reqStart := time.Now()
				resp, err := httpClient.Get(ctx, "/burst")
				duration := time.Since(reqStart)

				result := TestResult{
					WorkerID: requestID,
					Success:  err == nil,
					Duration: duration,
				}

				if err != nil {
					result.Error = err.Error()
				} else if resp.StatusCode != 200 {
					result.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
				}

				results <- result
			}(i)
		}

		wg.Wait()
		close(results)

		totalDuration := time.Since(start)

		// Analyze burst results
		successCount := 0
		var totalResponseTime time.Duration

		for result := range results {
			if result.Success {
				successCount++
				totalResponseTime += result.Duration
			}
		}

		successRate := float64(successCount) / float64(burstSize) * 100

		t.Logf("✅ Burst load test completed:")
		t.Logf("   Burst size: %d", burstSize)
		t.Logf("   Success rate: %.1f%%", successRate)
		t.Logf("   Total duration: %v", totalDuration)

		if successCount > 0 {
			avgDuration := totalResponseTime / time.Duration(successCount)
			t.Logf("   Average response time: %v", avgDuration)
		}

		// Assertions for burst test
		assert.Less(t, totalDuration, maxDuration,
			"Burst test should complete within time limit")
		assert.GreaterOrEqual(t, successRate, 70.0,
			"Even under burst load, success rate should be reasonable")
	})
}

// testMemoryUsage monitors memory usage during operations
func testMemoryUsage(t *testing.T) {
	// Create server that processes data
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return moderately large response
		response := map[string]interface{}{
			"data": make([]string, 100), // Array of 100 strings
			"metadata": map[string]string{
				"timestamp": time.Now().Format(time.RFC3339),
				"version":   "1.0",
			},
		}

		// Fill the data array
		for i := range response["data"].([]string) {
			response["data"].([]string)[i] = fmt.Sprintf("Data item %d with some content", i)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	ctx, err := createPerformanceCLIContextWithServer(server.URL)
	require.NoError(t, err)

	injector := di.Bootstrap(ctx)
	httpClient := di.GetHTTPClient(injector)

	// Get baseline memory stats
	runtime.GC()
	var baselineMem runtime.MemStats
	runtime.ReadMemStats(&baselineMem)

	// Perform memory-intensive operations
	const numOperations = 50
	for i := 0; i < numOperations; i++ {
		resp, err := httpClient.Get(context.Background(), "/data")
		if err != nil {
			t.Logf("⚠️ Request %d failed: %v", i, err)
			continue
		}

		if resp.StatusCode == 200 {
			// Process the response (parse JSON)
			var result map[string]interface{}
			if json.Unmarshal(resp.Body, &result) == nil {
				// Access the data to ensure it's loaded into memory
				if data, ok := result["data"].([]interface{}); ok {
					_ = len(data) // Access length
				}
			}
		}
	}

	// Check final memory usage
	runtime.GC()
	var finalMem runtime.MemStats
	runtime.ReadMemStats(&finalMem)

	memoryIncrease := float64(finalMem.Alloc-baselineMem.Alloc) / 1024 / 1024
	t.Logf("📊 Memory usage analysis:")
	t.Logf("   Baseline: %.2f MB", float64(baselineMem.Alloc)/1024/1024)
	t.Logf("   Final: %.2f MB", float64(finalMem.Alloc)/1024/1024)
	t.Logf("   Increase: %.2f MB", memoryIncrease)

	// Memory assertion
	assert.Less(t, memoryIncrease, 20.0,
		fmt.Sprintf("Memory increase should be reasonable, was %.2f MB", memoryIncrease))
}

// Helper types and functions

type TestResult struct {
	WorkerID int
	Success  bool
	Duration time.Duration
	Error    string
}

func createPerformanceCLIContextWithServer(serverURL string) (*cli.Context, error) {
	app := &cli.App{
		Name: "test",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "api-url",
				Value: serverURL,
			},
			&cli.StringFlag{
				Name:  "password",
				Value: "test",
			},
			&cli.BoolFlag{
				Name:  "verbose",
				Value: true,
			},
			&cli.IntFlag{
				Name:  "timeout",
				Value: 60, // Longer timeout for performance tests
			},
		},
	}

	flagSet := flag.NewFlagSet(app.Name, flag.ContinueOnError)
	flagSet.String("api-url", serverURL, "")
	flagSet.String("password", "test", "")
	flagSet.Bool("verbose", true, "")
	flagSet.Int("timeout", 60, "")

	args := []string{}
	err := flagSet.Parse(args)
	if err != nil {
		return nil, err
	}

	return cli.NewContext(app, flagSet, nil), nil
}
