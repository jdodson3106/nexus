package main

import "github.com/jdodson3106/nexus/cli/cmd/nexus"


func main() { 
   nexus.Execute() 
}

const helpText = `usage: nexus <command> [<args>...]
Nexus - An opinionated Web Framework in Go

commands:
    new     generates a new application
`

