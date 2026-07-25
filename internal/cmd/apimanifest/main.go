// Command apimanifest inventories every exported declaration in a pinned
// samber/lo checkout and records its GoForge compatibility status.
package main

import (
	"encoding/csv"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type symbol struct{ pkg, kind, name, file string }

var compatible = map[string]bool{
	"Entry": true, "Map": true, "Filter": true, "Reject": true,
	"FilterMap": true, "FlatMap": true, "Reduce": true, "ReduceRight": true,
	"ForEach": true, "Times": true, "Uniq": true, "UniqBy": true,
	"GroupBy": true, "GroupByMap": true, "KeyBy": true, "Associate": true,
	"Chunk": true, "Window": true, "Sliding": true, "PartitionBy": true,
	"Flatten": true, "Concat": true, "Interleave": true, "Reverse": true,
	"Drop": true, "DropRight": true, "Take": true, "Find": true,
	"FindIndexOf": true, "FindLastIndexOf": true, "First": true, "Last": true,
	"Contains": true, "Every": true, "Some": true, "None": true,
	"Nth":     true,
	"Without": true, "Union": true, "Intersect": true, "Difference": true,
	"Keys": true, "Values": true, "Entries": true, "FromEntries": true,
	"MapKeys": true, "MapValues": true,
}

func main() {
	if len(os.Args) != 3 {
		panic("usage: apimanifest UPSTREAM_ROOT OUTPUT.csv")
	}
	root, output := os.Args[1], os.Args[2]
	var symbols []symbol
	fset := token.NewFileSet()
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			name := entry.Name()
			if name == ".git" || name == "vendor" || (name == "internal" && filepath.Dir(path) == root) {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") || strings.HasSuffix(path, "_example.go") {
			return nil
		}
		file, parseErr := parser.ParseFile(fset, path, nil, 0)
		if parseErr != nil {
			return parseErr
		}
		rel, _ := filepath.Rel(root, filepath.Dir(path))
		if rel == "." {
			rel = "root"
		}
		add := func(kind, name string) {
			if ast.IsExported(name) {
				symbols = append(symbols, symbol{rel, kind, name, filepath.Base(path)})
			}
		}
		for _, declaration := range file.Decls {
			switch d := declaration.(type) {
			case *ast.FuncDecl:
				if d.Recv == nil {
					add("func", d.Name.Name)
				} else if ast.IsExported(d.Name.Name) {
					add("method", receiverName(d.Recv.List[0].Type)+"."+d.Name.Name)
				}
			case *ast.GenDecl:
				for _, specification := range d.Specs {
					switch spec := specification.(type) {
					case *ast.TypeSpec:
						add("type", spec.Name.Name)
					case *ast.ValueSpec:
						for _, name := range spec.Names {
							add(strings.ToLower(d.Tok.String()), name.Name)
						}
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	sort.Slice(symbols, func(i, j int) bool {
		a, b := symbols[i], symbols[j]
		if a.pkg != b.pkg {
			return a.pkg < b.pkg
		}
		if a.name != b.name {
			return a.name < b.name
		}
		if a.kind != b.kind {
			return a.kind < b.kind
		}
		return a.file < b.file
	})
	file, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	_ = writer.Write([]string{"package", "kind", "symbol", "source", "status", "destination_or_reason"})
	seen := map[string]bool{}
	for _, item := range symbols {
		key := item.pkg + "\x00" + item.kind + "\x00" + item.name
		if seen[key] {
			continue
		}
		seen[key] = true
		status, reason := "deferred", "outside compatibility tier 1"
		if item.pkg == "root" && compatible[item.name] {
			status, reason = "compatible", "goforge.dev/gplodash"
		}
		switch {
		case strings.HasPrefix(item.pkg, "exp/"):
			status, reason = "excluded-experimental", "upstream experimental/SIMD surface"
		case item.pkg == "mutable":
			status, reason = "deferred-mutable", "requires explicit aliasing API"
		case item.pkg == "parallel":
			status, reason = "deferred-parallel", "requires cancellation and ordering contract"
		case item.pkg == "it":
			status, reason = "deferred-iterator", "evaluate against Go iter.Seq and std algebra"
		}
		if err := writer.Write([]string{item.pkg, item.kind, item.name, item.file, status, reason}); err != nil {
			panic(err)
		}
	}
	if err := writer.Error(); err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stderr, "wrote %d unique exported declarations\n", len(seen))
}

func receiverName(expression ast.Expr) string {
	switch value := expression.(type) {
	case *ast.Ident:
		return value.Name
	case *ast.StarExpr:
		return receiverName(value.X)
	case *ast.IndexExpr:
		return receiverName(value.X)
	case *ast.IndexListExpr:
		return receiverName(value.X)
	}
	return "?"
}
