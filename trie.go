package iptrie

import (
	"inet.af/netaddr"
)

type TrieNode struct {
	children [2]*TrieNode
	prefix   netaddr.IPPrefix
	data     any
}

type CIDRTrie struct {
	rootV4 *TrieNode
	rootV6 *TrieNode
}

func NewCIDRTrie() *CIDRTrie {
	return &CIDRTrie{
		rootV4: &TrieNode{},
		rootV6: &TrieNode{},
	}
}

func (t *CIDRTrie) Insert(cidr string, data any) error {
	prefix, err := netaddr.ParseIPPrefix(cidr)
	if err != nil {
		return err
	}

	return t.InsertPrefix(prefix, data)
}

func (t *CIDRTrie) InsertPrefix(prefix netaddr.IPPrefix, data any) error {
	node := t.rootV4
	if prefix.IP().Is6() {
		node = t.rootV6
	}

	prefixLen := int(prefix.Bits())
	ip := prefix.IP()

	// Convert IP to bytes for bit manipulation
	var ipBytes []byte
	if ip.Is4() {
		arr := ip.As4()
		ipBytes = arr[:]
	} else {
		arr := ip.As16()
		ipBytes = arr[:]
	}

	for i := 0; i < prefixLen; i++ {
		byteIndex := i / 8
		bitIndex := 7 - (i % 8)
		bit := int((ipBytes[byteIndex] >> bitIndex) & 1)

		if node.children[bit] == nil {
			node.children[bit] = &TrieNode{}
		}
		node = node.children[bit]
	}
	node.prefix = prefix
	node.data = data
	return nil
}

func (t *CIDRTrie) SearchBest(ip netaddr.IP) (netaddr.IPPrefix, any) {
	node := t.rootV4

	// Handle IPv4-mapped IPv6 addresses
	if ip.Is6() {
		if ip.Is4in6() {
			// Convert IPv4-mapped IPv6 to IPv4 for lookup
			ip = ip.Unmap()
		} else {
			node = t.rootV6
		}
	}

	var bestMatch netaddr.IPPrefix
	var data any

	// Convert IP to bytes for bit manipulation
	var ipBytes []byte
	if ip.Is4() {
		arr := ip.As4()
		ipBytes = arr[:]
	} else {
		arr := ip.As16()
		ipBytes = arr[:]
	}

	// Directly calculate bits, avoiding slice allocation
	// Iterate up to 32 bits (IPv4) or 128 bits (IPv6)
	maxBits := len(ipBytes) * 8
	for i := 0; i < maxBits && node != nil; i++ {
		byteIndex := i / 8
		bitIndex := 7 - (i % 8)
		bit := (ipBytes[byteIndex] >> bitIndex) & 1

		node = node.children[bit]
		if node != nil && node.prefix.IsValid() {
			bestMatch = node.prefix
			data = node.data
		}
	}
	return bestMatch, data
}

func (t *CIDRTrie) SearchFast(ip netaddr.IP) (netaddr.IPPrefix, any) {
	node := t.rootV4

	// Handle IPv4-mapped IPv6 addresses
	if ip.Is6() {
		if ip.Is4in6() {
			// Convert IPv4-mapped IPv6 to IPv4 for lookup
			ip = ip.Unmap()
		} else {
			node = t.rootV6
		}
	}

	// Convert IP to bytes for bit manipulation
	var ipBytes []byte
	if ip.Is4() {
		arr := ip.As4()
		ipBytes = arr[:]
	} else {
		arr := ip.As16()
		ipBytes = arr[:]
	}

	// Directly calculate bits, avoiding slice allocation
	maxBits := len(ipBytes) * 8
	for i := 0; i < maxBits; i++ {
		byteIndex := i / 8
		bitIndex := 7 - (i % 8)
		bit := int((ipBytes[byteIndex] >> bitIndex) & 1)

		node = node.children[bit]
		if node == nil {
			break
		}
		if node.prefix.IsValid() {
			return node.prefix, node.data // 直接返回，无需临时变量
		}
	}
	return netaddr.IPPrefix{}, nil
}

func (t *CIDRTrie) IsBlank() bool {
	return t.rootV4.children[0] == nil &&
		t.rootV4.children[1] == nil &&
		t.rootV6.children[0] == nil &&
		t.rootV6.children[1] == nil
}

func (t *CIDRTrie) Contains(ip string) bool {
	netip, err := netaddr.ParseIP(ip)
	if err != nil {
		return false
	}
	prefix, _ := t.SearchFast(netip)
	return prefix.IsValid()
}

func (t *CIDRTrie) ContainsIP(netip netaddr.IP) bool {
	prefix, _ := t.SearchFast(netip)
	return prefix.IsValid()
}

// Legacy compatibility methods for net package
func (t *CIDRTrie) InsertIpNet(ipNet *netaddr.IPPrefix, data any) error {
	return t.InsertPrefix(*ipNet, data)
}
