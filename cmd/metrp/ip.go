package main

import (
	"fmt"
	"net"
	"strings"
)

// getPreferredIP returns a system IPv4 address that matches PreferredIPv4 or belongs to PreferredIPv4Prefix
func (app *application) getPreferredIP() (string, error) {
	if app.PreferredIPv4 != "" {
		ipv4Addr := net.ParseIP(app.PreferredIPv4)
		if ipv4Addr == nil {
			return "", fmt.Errorf("%q should be a valid IPv4 address", app.PreferredIPv4)
		}

		if strings.Contains(ipv4Addr.String(), ":") {
			return "", fmt.Errorf("%q should be an IPv4 address", app.PreferredIPv4)
		}

		return ipv4Addr.String(), nil
	}
	app.infoLog.Print("METRP_PREFERRED_IPV4 is not set, checking for METRP_PREFERRED_IPV4_PREFIX")

	if app.PreferredIPv4Prefix != "" {
		_, preferredPrefix, err := net.ParseCIDR(app.PreferredIPv4Prefix)
		if err != nil {
			return "", err
		}
		if strings.Contains(preferredPrefix.String(), ":") {
			return "", fmt.Errorf("%q should be an IPv4 prefix", app.PreferredIPv4Prefix)
		}

		// Get a list of unicast interface addresses
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			return "", err
		}

		// Go over IPv4 addresses and find one that belongs to the preferred network
		for _, addr := range addrs {
			if !strings.Contains(addr.String(), ":") {
				ipv4Addr, _, err := net.ParseCIDR(addr.String())
				if err != nil {
					return "", err
				}

				if preferredPrefix.Contains(ipv4Addr) {
					return ipv4Addr.String(), nil
				}
			}
		}

		return "", fmt.Errorf("There are no IPv4 interfaces that belong to METRP_PREFERRED_IPV4_PREFIX (%s)", app.PreferredIPv4Prefix)
	}

	app.infoLog.Print("METRP_PREFERRED_IPV4_PREFIX is not set either, so the app will listen on all interfaces")
	return "", nil
}
