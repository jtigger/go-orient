package main

import "testing"

func TestRoots(t *testing.T) {
	pkgs := NewPkgs()
	pkgs.AddDependence("a", "b")

	roots := pkgs.Roots()

	if len(roots) != 1 || roots[0].Name != "a" {
		t.Errorf("Expected exactly one root named '%s', but found: %+v", "a", roots)
	}
}

