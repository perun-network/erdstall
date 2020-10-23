// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"fmt"
	"text/tabwriter"
	"time"

	"github.com/montanaflynn/stats"
)

type Result struct {
	data  []float64
	begin time.Time
}

func Benchmark(n int, f func() error) (r Result, e error) {
	for i := 0; i < n; i++ {
		r.start()
		err := f()
		r.stop()
		if err != nil {
			return Result{}, err
		}
	}
	return
}

func (r *Result) start() {
	r.begin = time.Now()
}

func (r *Result) stop() {
	r.data = append(r.data, float64(time.Since(r.begin).Nanoseconds()))
}

func (r Result) String() string {
	functions := []func(stats.Float64Data) (float64, error){stats.Sum, stats.Min, stats.Max, stats.Median, stats.StdDevP}
	var str string
	values := make([]float64, len(functions))

	for i, f := range functions {
		values[i], _ = f(r.data)
		str += (time.Duration(values[i]) * time.Nanosecond).Round(time.Microsecond).String() + "\t"
	}

	var buff bytes.Buffer
	w := tabwriter.NewWriter(&buff, 0, 0, 3, ' ', tabwriter.Debug)

	freq := (float64(len(r.data)) / values[0]) * float64(time.Second.Nanoseconds())
	fmt.Fprintf(w, "N\ttx/s\tSum\tMin\tMax\tMedian\tStddev\t\n%d\t%.1f\t%s", len(r.data), freq, str)
	w.Flush()
	return buff.String()
}
