package sync_pool_performance

import (
	"runtime"
	"sync"
	"testing"

	"github.com/panjf2000/ants/v2"
)

var m map[string]interface{}

/*summary
BenchmarkVariablesPoolSyncPool-16                	    6103	    207113 ns/op	     204 B/op	       3 allocs/op
BenchmarkVariablesPoolSyncPoolWithChannel1k-16    	     238	   4904390 ns/op	    6355 B/op	      96 allocs/op
BenchmarkVariablesPoolSyncPoolWithChannel10k-16    	     202	   5709154 ns/op	   20323 B/op	     158 allocs/op
BenchmarkVariablesPoolSyncPoolWithChannel500-16    	     238	   4799105 ns/op	    5741 B/op	      94 allocs/op
 */

func runParallelWithAllocatedPool(b *testing.B, n int, f func()) {
	p, _ := ants.NewPool(n, ants.WithPreAlloc(true))

	cpus := runtime.NumCPU()
	p1, _ := ants.NewPool(cpus, ants.WithPreAlloc(true))

	wg := &sync.WaitGroup{}
	wg.Add(n * cpus)
	b.ReportAllocs()
	b.ResetTimer()
	//routine per n
	for i := 0; i < n; i++ {
		p.Submit(func() {
			// routine per cpu cores
			for k := 0; k < cpus; k++ {
				p1.Submit(func() {
					defer wg.Done()
					for j := 0; j < b.N; j++ {
						f()
					}
				})
			}
		})
	}
	wg.Wait()
}

//BenchmarkVariablesPoolSyncPool-16    	    6103	    207113 ns/op	     204 B/op	       3 allocs/op
func BenchmarkVariablesPoolSyncPool(b *testing.B) {
	vp := newVariablesMapPool(0, false, 5)
	runParallelWithAllocatedPool(b, 1000,  func() {
		v := vp.Get()
		v["key"] = "value"
		v["key2"] = "value2"
		m = v
		vp.Put(v)
	})
}

//BenchmarkVariablesPoolSyncPoolWithChannel1k-16    	     238	   4904390 ns/op	    6355 B/op	      96 allocs/op
func BenchmarkVariablesPoolSyncPoolWithChannel1k(b *testing.B) {
	vp := newVariablesMapPool(1000, true, 5)
	runParallelWithAllocatedPool(b, 1000,  func() {
		v := vp.Get()
		v["key"] = "value"
		v["key2"] = "value2"
		m = v
		vp.Put(v)
	})
}

//BenchmarkVariablesPoolSyncPoolWithChannel10k-16    	     202	   5709154 ns/op	   20323 B/op	     158 allocs/op
func BenchmarkVariablesPoolSyncPoolWithChannel10k(b *testing.B) {
	vp := newVariablesMapPool(10000, true, 5)
	runParallelWithAllocatedPool(b, 1000,  func() {
		v := vp.Get()
		v["key"] = "value"
		v["key2"] = "value2"
		m = v
		vp.Put(v)
	})
}

//BenchmarkVariablesPoolSyncPoolWithChannel500-16    	     238	   4799105 ns/op	    5741 B/op	      94 allocs/op
func BenchmarkVariablesPoolSyncPoolWithChannel500(b *testing.B) {
	vp := newVariablesMapPool(500, true, 5)
	runParallelWithAllocatedPool(b, 1000,  func() {
		v := vp.Get()
		v["key"] = "value"
		v["key2"] = "value2"
		m = v
		vp.Put(v)
	})
}