package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)


func main() {
	pkgs := NewPkgs()

	bytes, err := ioutil.ReadFile("deps.csv")
	if err != nil {
		panic(err)
	}
	r := csv.NewReader(strings.NewReader(string(bytes)))

	records, err := r.ReadAll()
	if err != nil {
		panic(err)
	}
	for _, record := range records {
		pkgs.AddDependence(record[0], record[1])
	}

	all := pkgs.All()
	sort.Slice(all, func(i, j int) bool {
		ideps := len(all[i].Dependants)
		jdeps := len(all[j].Dependants)
		if ideps != jdeps {
			return ideps < jdeps
		}
		ideps = len(all[i].Dependencies)
		jdeps = len(all[j].Dependencies)
		return ideps < jdeps
	})
	for _, pkg := range all {
		fmt.Printf("(%d) => %s => (%d)\n", len(pkg.Dependants), pkg.Name, len(pkg.Dependencies))
	}
	printed := NewPkgs()
	for _, root := range pkgs.Roots() {
		printDeps(root, printed)
	}
}

func printDeps(p *Pkg, printed *Pkgs) {
	if len(p.Dependencies) == 0 {
		printed.Add(p)
	}
	if printed.Has(p) {
		return
	}
	fmt.Printf("%s:\n", p.Name)
	deps := []*Pkg{}
	for pkg, _ := range p.Dependencies {
		deps = append(deps, pkg)
	}
	sort.Slice(deps, func(i, j int) bool {
		ideps := len(deps[i].Dependants)
		jdeps := len(deps[j].Dependants)
		if ideps != jdeps {
			return ideps < jdeps
		}
		return deps[i].Name < deps[j].Name
	})
	for _, dep := range deps {
		fmt.Printf("- %s\n", dep.Name)
	}
	for _, dep := range deps {
		printDeps(dep, printed)
		printed.Add(dep)
	}
}
