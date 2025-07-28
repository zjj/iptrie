package iptrie

import "net"

type TrieNode struct {
	children [2]*TrieNode
	cidr     *net.IPNet
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
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return err
	}

	if ipNet == nil {
		return nil
	}

	return t.InsertIpNet(ipNet, data)
}

func (t *CIDRTrie) InsertIpNet(ipNet *net.IPNet, data any) error {
	node := t.rootV4
	if ipNet.IP.To4() == nil {
		node = t.rootV6
	}

	prefix, _ := ipNet.Mask.Size()

	ip := ipNet.IP
	if ip.To4() != nil {
		ip = ip.To4()
	} else {
		ip = ip.To16()
	}

	for i := 0; i < prefix; i++ {
		byteIndex := i / 8
		bitIndex := 7 - (i % 8)
		bit := int((ip[byteIndex] >> bitIndex) & 1)

		if node.children[bit] == nil {
			node.children[bit] = &TrieNode{}
		}
		node = node.children[bit]
	}
	node.cidr = ipNet
	node.data = data
	return nil
}

func (t *CIDRTrie) SearchBest(ip net.IP) (*net.IPNet, any) {
	node := t.rootV4
	if ip.To4() == nil {
		node = t.rootV6
		// Use 16-byte representation for IPv6
		ip = ip.To16()
	} else {
		// Use 4-byte representation for IPv4
		ip = ip.To4()
	}

	if ip == nil {
		return nil, nil
	}

	var bestMatch *net.IPNet
	var data any

	maxBits := len(ip) * 8
	for i := 0; i < maxBits && node != nil; i++ {
		byteIndex := i / 8
		bitIndex := 7 - (i % 8)
		bit := int((ip[byteIndex] >> bitIndex) & 1)

		node = node.children[bit]
		if node != nil && node.cidr != nil {
			bestMatch = node.cidr
			data = node.data
		}
	}
	return bestMatch, data
}

func (t *CIDRTrie) SearchFast(ip net.IP) (*net.IPNet, any) {
	node := t.rootV4
	if ip.To4() == nil {
		node = t.rootV6
	}

	if ip.To4() != nil {
		ip = ip.To4()
	} else {
		ip = ip.To16()
	}

	maxBits := len(ip) * 8
	for i := 0; i < maxBits; i++ {
		byteIndex := i / 8
		bitIndex := 7 - (i % 8)
		bit := int((ip[byteIndex] >> bitIndex) & 1)

		node = node.children[bit]
		if node == nil {
			break
		}
		if node.cidr != nil {
			return node.cidr, node.data
		}
	}
	return nil, nil
}

func (t *CIDRTrie) IsBlank() bool {
	return t.rootV4.children[0] == nil &&
		t.rootV4.children[1] == nil &&
		t.rootV6.children[0] == nil &&
		t.rootV6.children[1] == nil
}

func (t *CIDRTrie) Contains(ip string) bool {
	netip := net.ParseIP(ip)
	if netip == nil {
		return false
	}
	n, _ := t.SearchFast(netip)
	return n != nil
}

func (t *CIDRTrie) ContainsIP(netip net.IP) bool {
	n, _ := t.SearchFast(netip)
	return n != nil
}
