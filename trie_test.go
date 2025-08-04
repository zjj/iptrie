package iptrie

import (
	"fmt"
	"testing"

	"inet.af/netaddr"
)

// setupTrie creates a trie with various CIDR blocks for testing
func setupTrie() *CIDRTrie {
	trie := NewCIDRTrie()

	// IPv4 ranges
	cidrs := []string{
		"192.168.0.0/16",
		"192.168.1.0/24",
		"192.168.1.128/25",
		"10.0.0.0/8",
		"172.16.0.0/12",
		"203.0.113.0/24",
		"198.51.100.0/24",
	}

	for i, cidr := range cidrs {
		trie.Insert(cidr, i)
	}

	// IPv6 ranges
	ipv6cidrs := []string{
		"2001:db8::/32",
		"2001:db8:1::/48",
		"2001:db8:1:1::/64",
		"fe80::/10",
		"::1/128",
	}

	for i, cidr := range ipv6cidrs {
		trie.Insert(cidr, i+100)
	}

	return trie
}

func BenchmarkSearchBest_IPv4_Hit(b *testing.B) {
	trie := setupTrie()
	ip, _ := netaddr.ParseIP("192.168.1.200")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.SearchBest(ip)
	}
}

func BenchmarkSearchBest_IPv4_Miss(b *testing.B) {
	trie := setupTrie()
	ip, _ := netaddr.ParseIP("8.8.8.8")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.SearchBest(ip)
	}
}

func BenchmarkSearchBest_IPv4_DeepMatch(b *testing.B) {
	trie := setupTrie()
	ip, _ := netaddr.ParseIP("192.168.1.150") // Should match the /25 subnet

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.SearchBest(ip)
	}
}

func BenchmarkSearchBest_IPv6_Hit(b *testing.B) {
	trie := setupTrie()
	ip, _ := netaddr.ParseIP("2001:db8:1:1::1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.SearchBest(ip)
	}
}

func BenchmarkSearchBest_IPv6_Miss(b *testing.B) {
	trie := setupTrie()
	ip, _ := netaddr.ParseIP("2001:db9::1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.SearchBest(ip)
	}
}

func BenchmarkSearchBest_LargeTrie(b *testing.B) {
	trie := NewCIDRTrie()

	// Create a larger trie with many subnets
	for i := 0; i < 256; i++ {
		for j := 0; j < 256; j += 16 {
			cidr := fmt.Sprintf("10.%d.%d.0/28", i, j)
			trie.Insert(cidr, i*256+j)
		}
	}

	ip, _ := netaddr.ParseIP("10.128.64.15")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.SearchBest(ip)
	}
}

func BenchmarkSearchBest_EmptyTrie(b *testing.B) {
	trie := NewCIDRTrie()
	ip, _ := netaddr.ParseIP("192.168.1.1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.SearchBest(ip)
	}
}

// Comparison benchmarks
func BenchmarkSearchFast_IPv4_Hit(b *testing.B) {
	trie := setupTrie()
	ip, _ := netaddr.ParseIP("192.168.1.200")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.SearchFast(ip)
	}
}

func BenchmarkSearchFast_IPv6_Hit(b *testing.B) {
	trie := setupTrie()
	ip, _ := netaddr.ParseIP("2001:db8:1:1::1")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.SearchFast(ip)
	}
}

// Test with different IP representations
func BenchmarkSearchBest_IPv4Mapped(b *testing.B) {
	trie := setupTrie()
	ip, _ := netaddr.ParseIP("::ffff:192.168.1.200")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.SearchBest(ip)
	}
}

// Memory allocation benchmark
func BenchmarkSearchBest_Allocs(b *testing.B) {
	trie := setupTrie()
	ip, _ := netaddr.ParseIP("192.168.1.200")

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.SearchBest(ip)
	}
}
