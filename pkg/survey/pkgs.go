package survey

type Pkg struct {
	Name         string
	Dependencies map[string]*Pkg
	Dependants   map[string]*Pkg
}

func (p *Pkg) String() string {
	return p.Name
}

func NewPkg(name string) *Pkg {
	return &Pkg{Name: name,
		Dependencies: make(map[string]*Pkg),
		Dependants:   make(map[string]*Pkg),
	}
}

type Pkgs struct {
	pkgs map[string]*Pkg
}

func NewPkgs() *Pkgs {
	return &Pkgs{pkgs: make(map[string]*Pkg)}
}

func (p *Pkgs) Has(pkg *Pkg) bool {
	_, has := p.pkgs[pkg.Name]
	return has
}

func (p *Pkgs) Get(pkgName string) Pkg {
	return *p.pkgs[pkgName]  // TODO: deep copy to ensure integrity of this Pkg
}

func (p *Pkgs) AddDependence(from, to string) {
	dependant, ok := p.pkgs[from]
	if !ok {
		dependant = NewPkg(from)
		p.pkgs[from] = dependant
	}
	dependency, ok := p.pkgs[to]
	if !ok {
		dependency = NewPkg(to)
		p.pkgs[to] = dependency
	}

	dependant.Dependencies[to] = dependency
	dependency.Dependants[from] = dependant
}

func (p *Pkgs) Roots() []*Pkg {
	roots := []*Pkg{}
	for _, pkg := range p.pkgs {
		if len(pkg.Dependants) == 0 {
			roots = append(roots, pkg)
		}
	}
	return roots
}

func (p *Pkgs) All() []*Pkg {
	all := []*Pkg{}
	for _, pkg := range p.pkgs {
		all = append(all, pkg)
	}
	return all
}
