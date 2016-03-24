package main

import (
  "github.com/docopt/docopt-go"
  "fmt"
)

func main() {
  usage := `gols.

Usage:
  gols --path=/this/path [--recursive] [--output=<text|json|yaml>]

Options:
  --path <path>       The path to the directory to list
  --recursive         List subdirectories recursively
  --output <format>   Output format [default: text]
  -h --help           Show program help`

  arguments, err := docopt.Parse(usage, nil, true, "", false)
  if (err) != nil {
    panic(err)
  }

  fmt.Println(arguments)
}
