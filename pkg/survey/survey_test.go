package survey

import (
	"testing"
)

func TestSurveyOver_CollectsPackagesWithinModule(t *testing.T) {
	basePath := "../.."
	survey, err := Of(basePath)
	if err != nil {
		t.Errorf("Of(\"%s\") failed: %s", basePath, err)
	}

	pkgs := survey.GetAllPackages()
	expectedPkgName := "pkg/survey"

    if !pkgs.Has(NewPkg(expectedPkgName)) {
    	t.Errorf("Expected to find package \"%s\"; packages are %v", expectedPkgName, pkgs.All())
	}
}

func TestSurvey_GetPackages_IncludesDependencies(t *testing.T) {
	basePath := "../.."
	survey, err := Of(basePath)
	if err != nil {
		t.Errorf("Of(\"%s\") failed: %s", basePath, err)
	}

	pkgs := survey.GetAllPackages()

	for _, root := range pkgs.Roots() {
		if len(root.Dependencies) < 1 {
			t.Errorf("expected roots to have at least one dependency; root \"%s\" had none", root.Name)
		}
	}
}