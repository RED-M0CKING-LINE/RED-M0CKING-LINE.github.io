// Go CIDR calculator for use in JavaScript
// Build with: GOOS=js GOARCH=wasm go build -o web/static/wasm/tool.wasm ./wasm/tool
// Registers a single global function, window.cidrCalc(string), which returns a structured object in web/static/js/tools.js
package main

import (
	"fmt"
	"math/big"
	"net/netip"
	"strings"
	"syscall/js" // This errors but its okay because this part is WASM
)

func main() {
	js.Global().Set("cidrCalc", js.FuncOf(cidrCalc))
	// Keep Go scheduler alive for callbacks
	select {}
}

func cidrCalc(_ js.Value, args []js.Value) any {
	if len(args) == 0 {
		return jsError("missing CIDR argument")
	}
	in := strings.TrimSpace(args[0].String())
	if in == "" {
		return jsError("empty input")
	}
	prefix, err := netip.ParsePrefix(in)
	if err != nil {
		return jsError("invalid CIDR: " + err.Error())
	}

	addr := prefix.Addr()
	bits := prefix.Bits()
	is4 := addr.Is4()
	totalBits := 32
	if !is4 {
		totalBits = 128
	}

	// Network address = mask of the input addr
	network := prefix.Masked().Addr()

	// Compute broadcast / last addr by setting all host bits to 1
	hostBits := totalBits - bits
	bcastBytes := append([]byte(nil), network.AsSlice()...)
	for i := len(bcastBytes) - 1; hostBits > 0 && i >= 0; i-- {
		take := 8
		if hostBits < 8 {
			take = hostBits
		}
		bcastBytes[i] |= byte((1 << take) - 1)
		hostBits -= take
	}
	broadcast, _ := netip.AddrFromSlice(bcastBytes)
	if is4 {
		broadcast = broadcast.Unmap()
	}

	// Masks
	maskBytes := make([]byte, totalBits/8)
	remaining := bits
	for i := range maskBytes {
		if remaining >= 8 {
			maskBytes[i] = 0xff
			remaining -= 8
		} else if remaining > 0 {
			maskBytes[i] = byte(0xff << (8 - remaining))
			remaining = 0
		}
	}
	wildBytes := make([]byte, len(maskBytes))
	for i, b := range maskBytes {
		wildBytes[i] = ^b
	}
	mask, _ := netip.AddrFromSlice(maskBytes)
	wild, _ := netip.AddrFromSlice(wildBytes)
	if is4 {
		mask = mask.Unmap()
		wild = wild.Unmap()
	}

	// Address counts as bigint to handle IPv6
	total := new(big.Int).Lsh(big.NewInt(1), uint(totalBits-bits))

	// Usable hosts: IPv4 with prefix < 31 reserves network + broadcast, /31 and /32 don't
	// IPv6 conventionally uses all
	usable := new(big.Int).Set(total)
	if is4 && bits < 31 {
		usable.Sub(usable, big.NewInt(2))
	}

	firstHost := network
	lastHost := broadcast
	if is4 && bits < 31 {
		firstHost = network.Next()
		lastHost = broadcast.Prev()
	}

	return map[string]any{
		"network":     network.String(),
		"broadcast":   broadcast.String(),
		"netmask":     mask.String(),
		"wildcard":    wild.String(),
		"firstHost":   firstHost.String(),
		"lastHost":    lastHost.String(),
		"prefix":      fmt.Sprintf("/%d", bits),
		"totalAddrs":  total.String(),
		"usableHosts": usable.String(),
	}
}

func jsError(msg string) map[string]any { return map[string]any{"error": msg} }
