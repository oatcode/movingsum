package movingsum_test

import (
	"fmt"
	"time"

	"github.com/oatcode/movingsum"
)

func ExampleNewMovingSum() {
	ms := movingsum.NewMovingSum(10)
	ms.Add(1)
	ms.Add(1)
	ms.Add(1)
	sum, count := ms.Get()
	avg := float64(sum) / float64(count)
	fmt.Printf("sum=%d avg=%.2f", sum, avg)
	// Output: sum=3 avg=1.00
}

func ExampleNewMovingSumByTime() {
	duration := 100 * time.Millisecond
	ms := movingsum.NewMovingSumByTime(duration)
	ms.Add(1)
	ms.Add(1)
	ms.Add(1)
	sum, count := ms.Get()
	avg := float64(sum) / float64(count)
	rate := float64(sum) / (float64(duration) / float64(time.Second))
	fmt.Printf("sum=%d avg=%.2f rate=%.2f", sum, avg, rate)
	// Output: sum=3 avg=1.00 rate=30.00
}
