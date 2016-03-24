package main

import (
  "github.com/docopt/docopt-go"
  "gopkg.in/yaml.v2"
  "encoding/json"
  "fmt"
  "io/ioutil"
  "os"
  "path/filepath"
  "time"
)

type NestedFile struct {
  ModifiedTime time.Time
  IsLink bool
  IsDir bool
  LinksTo string
  Size int64
  Name string
  Children []NestedFile
}

func getLinkInfo(f os.FileInfo, dir string) (bool, string) {
  if (f.Mode() & os.ModeSymlink != 0) {
    target, _ := os.Readlink(filepath.Join(dir, f.Name()))
    return true, target
  }
  return false, ""
}

func toNestedFile(f os.FileInfo, children []NestedFile, dir string) NestedFile {
  isLink, linkPath := getLinkInfo(f, dir)
  return NestedFile{
    f.ModTime(),
    isLink,
    f.IsDir(),
    linkPath,
    f.Size(),
    f.Name(),
    children,
  }
}

func dirReader(dir string, recursive bool) []NestedFile {
  files, _ := ioutil.ReadDir(dir)
  output := make([]NestedFile, 0, len(files))
  for _, f := range files {
    children := []NestedFile{}
    if f.IsDir() && recursive {
      children = dirReader(filepath.Join(dir, f.Name()), recursive)
    }
    output = append(output, toNestedFile(f, children, dir))
  }
  return output
}

func textOutputHelper(files []NestedFile, out *os.File, indent string) {
  for _, f := range files {
    fmt.Fprint(out, indent + f.Name)
    if f.IsDir {
      fmt.Fprintln(out, "/")
      textOutputHelper(f.Children, out, indent + "  ")
    } else if f.IsLink {
      fmt.Fprintln(out, "* (" + f.LinksTo + ")")
    } else {
      fmt.Fprintln(out, "")
    }
  }
}

func textOutput(files []NestedFile, path string, out *os.File) {
  fmt.Fprintln(out, path)
  textOutputHelper(files, out, "  ")
}

func jsonOutput(files []NestedFile, path string, out *os.File) {
  outJson, _ := json.MarshalIndent(files, "", "  ")
  fmt.Fprintln(out, string(outJson))
}

func yamlOutput(files []NestedFile, path string, out *os.File) {
  outYaml, _ := yaml.Marshal(files)
  fmt.Fprintln(out, string(outYaml))
}


func main() {
  usage := `gols.

Usage:
  gols --path=/this/path [--recursive] [--output=<text|json|yaml>]

Options:
  --path <path>       The path to the directory to list
  --recursive         List subdirectories recursively
  --output <format>   Output format [default: text]
  -h --help           Show program help`

  arguments, _ := docopt.Parse(usage, nil, true, "", false)

  format := arguments["--output"].(string)

  outputFormatters := map[string]func([]NestedFile, string, *os.File) {
    "json": jsonOutput,
    "text": textOutput,
    "yaml": yamlOutput,
  }

  formatter, ok := outputFormatters[format]
  if !ok {
    fmt.Println("Invalid output format provided.")
    os.Exit(1)
  }

  path := arguments["--path"].(string)
  absPath, _ := filepath.Abs(path)
  recursive := arguments["--recursive"].(bool)
  formatter(dirReader(absPath, recursive), absPath, os.Stdout)
}
