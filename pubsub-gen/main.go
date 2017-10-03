package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"code.cloudfoundry.org/go-pubsub/pubsub-gen/internal/generator"
	"code.cloudfoundry.org/go-pubsub/pubsub-gen/internal/inspector"
)

func main() {
	structPath := flag.String("struct-name", "", "The name of the struct create a traverser for")
	packageName := flag.String("package", "", "The package name of the generated code")
	traverserName := flag.String("traverser", "", "The name of the generated traverser")
	output := flag.String("output", "", "The path to output the generated file")
	isPtr := flag.Bool("pointer", false, "Will the struct be a pointer when being published?")
	includePkgName := flag.Bool("include-pkg-name", false, "Prefix the struct type with the package name?")
	interfaces := flag.String("interfaces", "{}", "A map (map[string][]string encoded in JSON) mapping interface types to implementing structs")
	slices := flag.String("slices", "{}", "A map (map[string][]string encoded in JSON) mapping types to field names for slices")
	subStructs := flag.String("sub-structs", "{}", "A map (map[string]string encoded in JSON) mapping names to package locations")
	imports := flag.String("imports", "", "A comma separated list of imports required in the generated file")
	blacklist := flag.String("blacklist-fields", "", `A comma separated list of struct name and field
	combos to not include (e.g., mystruct.myfield,otherthing.otherfield).
	A wildcard (*) can be provided for the struct name (e.g., *.fieldname).`)

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

	ms := make(map[string]string)
	if err := json.Unmarshal([]byte(*subStructs), &ms); err != nil {
		log.Fatalf("Invalid sub-structs (%s): %s", *subStructs, err)
	}

	sliceM := make(map[string]string)
	if err := json.Unmarshal([]byte(*slices), &sliceM); err != nil {
		log.Fatalf("Invalid slices (%s): %s", *slices, err)
	}

	importList := strings.Split(*imports, ",")

	structName := (*structPath)[idx+1:]

	var pkgName string
	if *includePkgName {
		idx2 := strings.LastIndex(filepath.ToSlash(*structPath), "/")
		pkgName = filepath.ToSlash(*structPath)[idx2+1:idx] + "."
	}

	fieldBlacklist := buildBlacklist(*blacklist)

	sf := inspector.NewStructFetcher(fieldBlacklist, ms, sliceM)
	pp := inspector.NewPackageParser(sf)

	mm := make(map[string]inspector.Struct)
	for fullName, path := range ms {
		splitName := strings.SplitN(fullName, ".", 2)
		var name, pkg string
		_ = name
		if len(splitName) != 2 {
			name = fullName
		} else {
			pkg = splitName[0]
			name = splitName[1]
		}

		m, err := pp.Parse(path, gopath)
		if err != nil {
			log.Fatal(err)
		}

		for k, v := range m {
			v.Name = fmt.Sprintf("%s.%s", pkg, k)
			if strings.HasSuffix(fullName, v.Name) {
				mm[v.Name] = v
			} else {
				mm[k] = v
			}
		}
	}

	m, err := pp.Parse((*structPath)[:idx], gopath)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range m {
		mm[k] = v
	}

	linker := inspector.NewLinker()
	linker.Link(mm, mi)

	g := generator.NewTraverserGenerator(generator.CodeWriter{})
	src, err := g.Generate(
		mm,
		*packageName,
		*traverserName,
		structName,
		*isPtr,
		pkgName,
		importList,
	)
	if err != nil {
		log.Fatal(err)
	}

	pg := generator.NewPathGenerator()
	src, err = pg.Generate(src, mm, *traverserName, structName)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(*output, []byte(src), 420)
	if err != nil {
		log.Fatal(err)
	}
}

func buildBlacklist(bl string) map[string][]string {
	if len(bl) == 0 {
		return nil
	}

	m := make(map[string][]string)
	for _, s := range strings.Split(bl, ",") {
		x := strings.Split(s, ".")
		if len(x) != 2 {
			log.Fatalf("'%s' is not in the proper format (structname.fieldname)", x)
		}

		m[x[0]] = append(m[x[0]], x[1])
	}
	return m
}
