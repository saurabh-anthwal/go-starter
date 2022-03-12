/*
Copyright © 2021 Rewire Group, Inc. All rights reserved.

Proprietary and confidential.

Unauthorized copying or use of this file, in any medium or form,
is strictly prohibited.
*/

package config

import (
	"flag"
	"fmt"
)

// Config holds configuration for the Server.
type Config struct {
	host string
	port string
}

// Instance holds the parsed flags.
var Flags = &Config{}

// Flags define the command line flags for the controller
func DefineFlags() {
	flag.StringVar(&Flags.host, "host", "127.0.0.1", "Host on which to run the HTTP server")
	flag.StringVar(&Flags.port, "port", "8080", "Port for HTTP server")
}

/// Define getters for each of the members to ensure immutability.

func (s *Config) GetHost() string {
	return s.host
}

func (s *Config) GetPort() string {
	return s.port
}

func (s *Config) GetHostPort() string {
	return fmt.Sprintf("%s:%s", s.host, s.port)
}

// Validate does necessary Flags validation and parsing.
// Its ok to fatal here.
func (s *Config) Validate() {

}
