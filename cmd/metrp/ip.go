package main

import (
	"fmt"
	"net"
	"net/netip"
)

// getPreferredIP returns a system IPv4 address that matches PreferredIPv4 or belongs to PreferredIPv4Prefix
func (app *application) getPreferredIP() (string, error) {
	if app.PreferredIPv4 != "" {
		ip, err := netip.ParseAddr(app.PreferredIPv4)
		if err != nil {
			return "", fmt.Errorf("%q should be a valid IPv4 address", app.PreferredIPv4)
		}

		if !ip.Is4() {
			return "", fmt.Errorf("%q should be an IPv4 address", app.PreferredIPv4)
		}

		return ip.String(), nil
	}
	app.infoLog.Print("METRP_PREFERRED_IPV4 is not set, checking for METRP_PREFERRED_IPV4_PREFIX")

	if app.PreferredIPv4Prefix != "" {
		preferredPrefix, err := netip.ParsePrefix(app.PreferredIPv4Prefix)
		if err != nil {
			return "", err
		}

		if !preferredPrefix.Addr().Is4() {
			return "", fmt.Errorf("%q should be an IPv4 prefix", app.PreferredIPv4Prefix)
		}

		// Get a list of unicast interface addresses
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			return "", err
		}

		// Go over IPv4 addresses and find one that belongs to the preferred network
		for _, addr := range addrs {
			prefix, err := netip.ParsePrefix(addr.String())
			// Should never happen
			if err != nil {
				app.errorLog.Print(err)
				continue
			}

			if preferredPrefix.Contains(prefix.Addr()) {
				return prefix.Addr().String(), nil
			}
		}

		return "", fmt.Errorf("there are no IPv4 interfaces that belong to METRP_PREFERRED_IPV4_PREFIX (%s)", app.PreferredIPv4Prefix)
	}

	app.infoLog.Print("METRP_PREFERRED_IPV4_PREFIX is not set either, so the app will listen on all interfaces")
	return "", nil
}
