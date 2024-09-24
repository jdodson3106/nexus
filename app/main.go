package main

import "github.com/jdodson3106/nexus"

func main() {

	nx, err := nexus.InitNexus()

	if err != nil {
		panic(err)
	}

	if err := nx.Run(); err != nil {
		panic(err)
	}
}
