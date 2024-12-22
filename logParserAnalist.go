/*
Log Parser Analytics

You're building a log analysis tool for a web service. You have log entries in the following format:
timestamp|ip_address|user_id|action|resource|status_code

Example log line:
2024-01-15T14:22:33Z|192.168.1.1|user123|GET|/api/products|200

Write a program that:
1. Reads log entries from a file or stdin (one per line)
2. Processes them concurrently using goroutines
3. Calculates the following metrics:
   - Top 3 most accessed resources
   - Number of unique users per hour
   - Success rate (percentage of 2xx status codes)
   - Average response time per endpoint (if status code is 2xx)
4. Can handle malformed lines gracefully
5. Outputs results in JSON format
6. Has a timeout for processing (e.g., 5 seconds)

Additional requirements:
- Use channels for communication between goroutines
- Implement proper error handling
- Use context for timeout
- Use sync.WaitGroup for coordination
- Use time.Time for timestamp parsing
- Use maps for aggregation
- Use mutexes for thread-safe operations

Sample Input:
2024-01-15T14:22:33Z|192.168.1.1|user123|GET|/api/products|200
2024-01-15T14:23:33Z|192.168.1.2|user456|GET|/api/orders|500
2024-01-15T14:24:33Z|192.168.1.1|user123|GET|/api/products|200
2024-01-15T15:22:33Z|192.168.1.3|user789|GET|/api/users|404
*/

package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type LogEntry struct {
	Timestamp  time.Time
	IP         string
	UserID     string
	Action     string
	Resource   string
	StatusCode int
}

type Metrics struct {
	TopResources     []ResourceCount    `json:"top_resources"`
	UniqueUsersPerHr map[string]int     `json:"users_per_hour"`
	SuccessRate      float64            `json:"success_rate"`
	AverageRespTime  map[string]float64 `json:"avg_response_time"`
}

type ResourceCount struct {
	Resource string `json:"resource"`
	Count    int    `json:"count"`
}

type MetricsCollector struct {
	mu              sync.Mutex
	resourceCounts  map[string]int
	uniqueUsers     map[string]map[string]bool
	totalRequests   int
	successRequests int
	responseTimes   map[string][]float64
}

func parseLogLine(line string) (LogEntry, error) {
	parts := strings.Split(line, "|")
	if len(parts) != 6 {
		return LogEntry{}, fmt.Errorf("invalid log format: %s", line)
	}

	timestamp, err := time.Parse(time.RFC3339, parts[0])
	if err != nil {
		return LogEntry{}, fmt.Errorf("invalid timestamp: %v", err)
	}

	statusCode, err := strconv.Atoi(parts[5])
	if err != nil {
		return LogEntry{}, fmt.Errorf("invalid status code: %v", err)
	}

	return LogEntry{
		Timestamp:  timestamp,
		IP:         parts[1],
		UserID:     parts[2],
		Action:     parts[3],
		Resource:   parts[4],
		StatusCode: statusCode,
	}, nil
}

func newMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		resourceCounts: make(map[string]int),
		uniqueUsers:    make(map[string]map[string]bool),
		responseTimes:  make(map[string][]float64),
	}
}

func processLogs(ctx context.Context, input <-chan string) (*Metrics, error) {
	collector := newMetricsCollector()
	var wg sync.WaitGroup

	// Process logs concurrently
	for i := 0; i < 5; i++ { // Use 5 worker goroutines
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case line, ok := <-input:
					if !ok {
						return
					}
					entry, err := parseLogLine(line)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error parsing log: %v\n", err)
						continue
					}
					collector.processEntry(entry)
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	wg.Wait()
	return collector.computeMetrics()
}

func (mc *MetricsCollector) processEntry(entry LogEntry) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	// Update resource counts
	mc.resourceCounts[entry.Resource]++

	// Update unique users per hour
	hour := entry.Timestamp.Format("2006-01-02-15")
	if mc.uniqueUsers[hour] == nil {
		mc.uniqueUsers[hour] = make(map[string]bool)
	}
	mc.uniqueUsers[hour][entry.UserID] = true

	// Update success rate
	mc.totalRequests++
	if entry.StatusCode >= 200 && entry.StatusCode < 300 {
		mc.successRequests++
	}
}

func (mc *MetricsCollector) computeMetrics() (*Metrics, error) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	// Compute top resources
	var resources []ResourceCount
	for resource, count := range mc.resourceCounts {
		resources = append(resources, ResourceCount{Resource: resource, Count: count})
	}
	sort.Slice(resources, func(i, j int) bool {
		return resources[i].Count > resources[j].Count
	})

	// Get top 3 or less
	if len(resources) > 3 {
		resources = resources[:3]
	}

	// Compute unique users per hour
	uniqueUsersPerHr := make(map[string]int)
	for hour, users := range mc.uniqueUsers {
		uniqueUsersPerHr[hour] = len(users)
	}

	// Compute success rate
	successRate := 0.0
	if mc.totalRequests > 0 {
		successRate = float64(mc.successRequests) / float64(mc.totalRequests) * 100
	}

	return &Metrics{
		TopResources:     resources,
		UniqueUsersPerHr: uniqueUsersPerHr,
		SuccessRate:      successRate,
		AverageRespTime:  make(map[string]float64), // Implement if response time is available in logs
	}, nil
}

func logParserAnalist() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create channel for log lines
	logChan := make(chan string, 100)

	// Start reading input in a separate goroutine
	go func() {
		defer close(logChan)
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			case logChan <- scanner.Text():
			}
		}
	}()

	// Process logs
	metrics, err := processLogs(ctx, logChan)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing logs: %v\n", err)
		os.Exit(1)
	}

	// Output results as JSON
	output, err := json.MarshalIndent(metrics, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding metrics: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(output))
}
