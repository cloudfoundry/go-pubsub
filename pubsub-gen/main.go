package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/apoydence/pubsub/pubsub-gen/internal/generator"
	"github.com/apoydence/pubsub/pubsub-gen/internal/inspector"
)

func main() {
	structPath := flag.String("struct-name", "", "The name of the struct create a traverser for")
	packageName := flag.String("package", "", "The package name of the generated code")
	traverserName := flag.String("traverser", "", "The name of the generated traverser")
	output := flag.String("output", "", "The path to output the generated file")
	isPtr := flag.Bool("pointer", false, "Will the struct be a pointer when being published?")
	interfaces := flag.String("interfaces", "{}", "A map (map[string][]string encoded in JSON) mapping interface types to implementing structs")

	flag.Parse()
	gopath := os.Getenv("GOPATH")

	if gopath == "" {
		log.Fatal("GOPATH is empty")
	}

	if *structPath == "" {
		log.Fatal("struct-name is required")
	}

	if *packageName == "" {
		log.Fatal("package is required")
	}

	if *traverserName == "" {
		log.Fatal("traverser is required")
	}

	if *output == "" {
		log.Fatal("output is required")
	}

	idx := strings.LastIndex(filepath.ToSlash(*structPath), ".")
	if idx < 0 {
		log.Fatalf("Invalid struct name: %s", *structPath)
	}

	mi := make(map[string][]string)
	if err := json.Unmarshal([]byte(*interfaces), &mi); err != nil {
		log.Fatalf("Invalid interfaces (%s): %s", *interfaces, err)
	}

	structName := (*structPath)[idx+1:]

	sf := inspector.NewStructFetcher()
	pp := inspector.NewPackageParser(sf)
	m, err := pp.Parse((*structPath)[:idx], gopath)
	if err != nil {
		log.Fatal(err)
	}

	linker := inspector.NewLinker()
	linker.Link(m, mi)

	g := generator.New(generator.CodeWriter{})
	src, err := g.Generate(m, *packageName, *traverserName, structName, *isPtr)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(*output, []byte(src), 420)
	if err != nil {
		log.Fatal(err)
	}
}
