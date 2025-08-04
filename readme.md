```
go test -bench=BenchmarkSearchBest -benchmem
goos: linux
goarch: amd64
pkg: github.com/zjj/iptrie
cpu: AMD Ryzen 7 PRO 4750U with Radeon Graphics
BenchmarkSearchBest_IPv4_Hit-16          	26202474	        44.64 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchBest_IPv4_Miss-16         	83709484	        13.16 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchBest_IPv4_DeepMatch-16    	26962598	        44.70 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchBest_IPv6_Hit-16          	11062160	       108.2 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchBest_IPv6_Miss-16         	21149138	        56.42 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchBest_LargeTrie-16         	22657850	        53.68 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchBest_EmptyTrie-16         	262643415	         4.553 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchBest_IPv4Mapped-16        	26265601	        45.26 ns/op	       0 B/op	       0 allocs/op
BenchmarkSearchBest_Allocs-16            	26688318	        44.85 ns/op	       0 B/op	       0 allocs/op
```
