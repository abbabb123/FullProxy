package ArgumentParser

import (
	"flag"
	"os"
	"time"
)

func ParseForwardToSocks5Arguments() (*string, *string, *string, *string, *string, *string, *string, *string, *int, *time.Duration) {
	protocolFlagSet := flag.NewFlagSet("port_forward-socks5", flag.ExitOnError)
	bindHost := protocolFlagSet.String("bind-host", "0.0.0.0", "Host to listen on.")
	bindPort := protocolFlagSet.String("bind-port", "8080", "Port to listen on.")
	socks5Host := protocolFlagSet.String("socks5-host", "127.0.0.1", "SOCKS5 server host to use")
	socks5Port := protocolFlagSet.String("socks5-port", "1080", "SOCKS5 server port to use")
	username := protocolFlagSet.String("socks5-username", "", "Username for the SOCKS5 server; leave empty for no AUTH")
	password := protocolFlagSet.String("socks5-password", "", "Password for the SOCKS5 server; leave empty for no AUTH")
	targetHost := protocolFlagSet.String("target-host", "", "Host of the target host that is accessible by the SOCKS5 proxy")
	targetPort := protocolFlagSet.String("target-port", "", "Port of the target host that is accessible by the SOCKS5 proxy")
	tries := protocolFlagSet.Int("tries", 5, "The number of re-tries that will maintain the connection between target and client (default is 5 tries)")
	rawTimeout := protocolFlagSet.Int("timeout", 10, "The number of second before re-trying the connection between target and client (default is 10 seconds)")
	_ = protocolFlagSet.Parse(os.Args[3:])
	timeout := time.Duration(*rawTimeout) * time.Second
	return bindHost, bindPort, socks5Host, socks5Port, username, password, targetHost, targetPort, tries, &timeout
}