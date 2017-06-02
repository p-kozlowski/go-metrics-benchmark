package go_metrics_benchmark

import (
	"fmt"
	"github.com/rcrowley/go-metrics"
	"sync"
	"testing"
	"time"
)

func measureRequestTime(timeStarted time.Time, requestType string, schemaId string) {
	m := metrics.GetOrRegisterTimer(fmt.Sprintf("req.%s.%s", schemaId, requestType), metrics.DefaultRegistry)
	m.UpdateSince(timeStarted)
}

func BenchmarkMetrics(b *testing.B) {
	for n := 0; n < b.N; n++ {
		measureRequestTime(time.Now(), "update", "dummyId")
	}
}

func measureInParallel(group *sync.WaitGroup) {
	measureRequestTime(time.Now(), "update", "firstId")
	measureRequestTime(time.Now(), "update", "secondId")
	group.Done()
}

func BenchmarkMetricsLockContention(b *testing.B) {
	group := sync.WaitGroup{}
	group.Add(b.N)

	for n := 0; n < b.N; n++ {
		go measureInParallel(&group)
	}

	group.Wait()
}
