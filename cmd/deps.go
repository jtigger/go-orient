package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/jtigger/go-orient/pkg/survey"
)


func main() {
	path := "."
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	survey, err := survey.Of(path)
	if err != nil {
		panic(fmt.Errorf("survey over \"%s\" failed: %s", path, err))
	}
	pkgs := survey.GetCodePackages()

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
	printed := map[string]bool{}
	for _, root := range pkgs.Roots() {
		printDeps(root, printed)
	}
}


func printDeps(p *survey.Pkg, printed map[string]bool) {
	if len(p.Dependencies) == 0 {
		printed[p.Name]=true
	}
	if _, has := printed[p.Name]; has {
		return
	}
	fmt.Printf("%s:\n", p.Name)
	var deps []*survey.Pkg
	for _, pkg := range p.Dependencies {
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
		printed[p.Name]=true
	}
}
