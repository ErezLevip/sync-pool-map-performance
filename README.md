# Sync Pool Map Performance Benchmark
    goos: darwin
    goarch: amd64
    pkg: sync-pool
    cpu: Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz

    BenchmarkVariablesPoolSyncPool-16                	    6103	    207113 ns/op	     204 B/op	       3 allocs/op
    BenchmarkVariablesPoolSyncPoolWithChannel1k-16    	     238	   4904390 ns/op	    6355 B/op	      96 allocs/op
    BenchmarkVariablesPoolSyncPoolWithChannel10k-16    	     202	   5709154 ns/op	   20323 B/op	     158 allocs/op
    BenchmarkVariablesPoolSyncPoolWithChannel500-16    	     238	   4799105 ns/op	    5741 B/op	      94 allocs/op