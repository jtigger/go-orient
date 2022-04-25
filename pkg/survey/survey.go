package survey

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/jtigger/go-orient/pkg/gomod"
)

type Survey struct {
	mod     *gomod.GoMod
	astPkgs map[string]*ast.Package
	pkgs    *Pkgs
}

// Of collects structural data about the code in the Go Module rooted at `relBasePath`.
//   `relBasePath` can be a relative path.
func Of(relBasePath string) (*Survey, error) {
	basePath, err := filepath.Abs(relBasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve \"%s\" to an absolute path: %s", relBasePath, err)
	}
	mod, err := gomod.GetModule(basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to obtain Go Module info from module at %s: %s", basePath, err)
	}

	fileSet := &token.FileSet{}
	pkgs := map[string]*ast.Package{}

	err = filepath.Walk(basePath, func(pkgPath string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}
		if info.Name() == ".git" {
			return filepath.SkipDir
		}

		pkgInDir, err := parser.ParseDir(fileSet, pkgPath, nil, parser.AllErrors)
		for _, pkg := range pkgInDir {
			pkgs[qualifiedPkgName(pkg.Name, pkgPath, basePath)] = pkg
		}
		if err != nil {
			return fmt.Errorf("surveying %s: %s", basePath, err)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &Survey{mod: mod, astPkgs: pkgs}, nil
}

// qualifiedPkgName calculates the path-qualified name of the package at `pkgPath`.
//    `pkgPath` is assumed to be `basePath` plus the path-parts to the directory
//    containing the package.
func qualifiedPkgName(pkgName, pkgPath, basePath string) string {
	var name string
	switch {
	case pkgName == "main":
		name = pkgName
	case filepath.Base(pkgPath) == pkgName:
		name = pkgPath[len(basePath)+1:]
	default:
		name = filepath.Dir(pkgPath[len(basePath)+1:]) + pkgName
	}
	return name
}

func (s *Survey) GetAllPackages() Pkgs {
	if s.pkgs == nil {
		s.pkgs = NewPkgs()
		for qualPkgName, astPkg := range s.astPkgs {
			qualDepPkgNames := s.dependenciesFor(qualPkgName, astPkg)
			for _, qualDepPkgName := range qualDepPkgNames {
				s.pkgs.AddDependence(qualPkgName, qualDepPkgName)
			}
		}
	}
	return *s.pkgs
}

func (s *Survey) GetCodePackages() Pkgs {
	if s.pkgs == nil {
		s.pkgs = NewPkgs()
		for qualPkgName, astPkg := range s.astPkgs {
			qualDepPkgNames := s.dependenciesFor(qualPkgName, astPkg)
			for _, qualDepPkgName := range qualDepPkgNames {
				if strings.HasSuffix(qualPkgName, "_test") ||
					strings.HasSuffix(qualDepPkgName, "_test") {
					continue
				}
				s.pkgs.AddDependence(qualPkgName, qualDepPkgName)
			}
		}
	}
	return *s.pkgs
}

func (s *Survey) dependenciesFor(qualPkgName string, astPkg *ast.Package) []string {
	dependencies := []string{}
	for _, file := range astPkg.Files {
		for _, importSpec := range file.Imports {
			depPkgPath := strings.Trim(importSpec.Path.Value, "\"")
			if !strings.HasPrefix(depPkgPath, s.mod.Module.Path) {
				continue
			}
			dependencies = append(dependencies, depPkgPath[len(s.mod.Module.Path)+1:])
		}
	}
	return dependencies
}

type visitor struct {
	qualPkgName string
	invocations []Invocation
}

type Invocation struct {
	QualPkgName string
	TypeName    string
	FunName     string
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	if node == nil {
		return nil
	}
	switch n := node.(type) {
	case *ast.CallExpr:
		var pkgName string
		var name string
		fmt.Printf("invokes function \"%s()\"\n", n.Fun)
		switch funName := n.Fun.(type) {
		case *ast.Ident:
			pkgName = ""
			name = funName.Name
		case *ast.SelectorExpr:
			switch x := funName.X.(type) {
			case *ast.Ident:
				pkgName = x.Name
			}
			name = funName.Sel.Name
		}
		v.invocations = append(v.invocations, Invocation{
			QualPkgName: pkgName,
			TypeName:    "",
			FunName:     name,
		})
	}
	fmt.Printf("<%T> %v\n", node, node)
	return v
}

func (v *visitor) Invocations() []Invocation {
	return v.invocations
}

func (s *Survey) GetFunctions() []Invocation {
	var invocations []Invocation
	for qualPkgName, astPkg := range s.astPkgs {
		funNamesPlucker := &visitor{
			qualPkgName: qualPkgName,
		}
		for _, file := range astPkg.Files {
			ast.Walk(funNamesPlucker, file)
		}
		invocations = append(invocations, funNamesPlucker.Invocations()...)
	}
	return invocations
}
