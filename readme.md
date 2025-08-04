```
go test -bench=BenchmarkSearchBest -benchmem
goos: linux
goarch: amd64
pkg: github.com/zjj/iptrie
cpu: AMD Ryzen 7 PRO 4750U with Radeon Graphics
BenchmarkSearchBest_IPv4_Hit-16                 21462196                53.55 ns/op            0 B/op          0 allocs/op
BenchmarkSearchBest_IPv4_Miss-16                57089424                20.90 ns/op            0 B/op          0 allocs/op
BenchmarkSearchBest_IPv4_DeepMatch-16           22798498                52.55 ns/op            0 B/op          0 allocs/op
BenchmarkSearchBest_IPv6_Hit-16                  9710242               114.4 ns/op             0 B/op          0 allocs/op
BenchmarkSearchBest_IPv6_Miss-16                20364656                58.28 ns/op            0 B/op          0 allocs/op
BenchmarkSearchBest_LargeTrie-16                19756615                58.65 ns/op            0 B/op          0 allocs/op
BenchmarkSearchBest_EmptyTrie-16                100000000               11.45 ns/op            0 B/op          0 allocs/op
BenchmarkSearchBest_IPv4Mapped-16               22593169                52.85 ns/op            0 B/op          0 allocs/op
BenchmarkSearchBest_Allocs-16                   22488368                52.47 ns/op            0 B/op          0 allocs/op
```