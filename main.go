/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"embed"

	"github.com/nahtann/trancome/cmd"
)

//go:embed migrations/shared/*sql
var migrations embed.FS

func main() {
	cmd.Execute(migrations)
}
