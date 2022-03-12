package main

import "github.com/saurabh-anthwal/dummy/pkg/config"
import "github.com/saurabh-anthwal/dummy/server"

func main()  {

	config.DefineFlags()
	server.StartServer(config.Flags)
}

