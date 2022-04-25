# go-orient

Get acquainted with a Go codebase.

Today, `go-orient` reports on _internal_ dependencies between packages.
That is, how: the packages of this module relate

## Usage

Install the binary

```console
go install github.com/jtigger/go-orient@latest
```
_(As of now, this tool is not versioned.)_

In the root of your Go project, run the tool:

```console
$ cd <ROOT-OF-GO-PROJECT>
$ go-orient
```

## Report Format

There are two parts to the "report".

The first port summarizes dependency cardinalities:

```
(# of dependents) => (package name) => (# of dependencies)
```

The second half details direct dependencies of each package.

```
- (package name)
  - (dependency 1 package name)
  - (dependency 2 package name)
  - ...
```

For example, a report run against `vmware-tanzu/carvel-ytt` looked like this:

```
(0) => main => (3)
(0) => pkg/cmd => (8)
(1) => pkg/website => (0)
(1) => pkg/yamlmeta/internal/yaml.v2 => (0)
(1) => pkg/spell => (0)
(1) => pkg/yamlfmt => (1)
(1) => pkg/yttlibraryext => (1)
(1) => pkg/yttlibraryext/toml => (4)
(1) => pkg/workspace => (13)
(2) => pkg/version => (0)
(2) => pkg/texttemplate => (2)
(2) => pkg/yttlibrary => (5)
(2) => pkg/workspace/datavalues => (6)
(2) => pkg/yamltemplate => (6)
(2) => pkg/cmd/template => (10)
(3) => pkg/workspace/ref => (2)
(3) => pkg/yttlibrary/overlay => (5)
(3) => pkg/schema => (5)
(4) => pkg/files => (1)
(5) => pkg/orderedmap => (0)
(5) => pkg/cmd/ui => (0)
(8) => pkg/template => (2)
(9) => pkg/filepos => (0)
(9) => pkg/template/core => (1)
(10) => pkg/yamlmeta => (3)
...
```

Where:
- `main` package has not dependents and 3 dependencies
- `pkg/yamlmeta` package has 10 dependents and 3 dependencies.
- `pkg/orderedmap` package has 5 dependents and no dependencies.


and the report continues...
```
pkg/cmd:
- pkg/website
- pkg/yamlfmt
- pkg/yttlibraryext
- pkg/cmd/template
- pkg/version
- pkg/files
- pkg/cmd/ui
- pkg/yamlmeta
```

Where:
- the `pkg/cmd` package depends on (among the others listed), `pkg/cmd/template` and `pkg/yamlmeta` .

and so forth for each package in the project:

```
pkg/yamlfmt:
- pkg/yamlmeta
pkg/yamlmeta:
- pkg/yamlmeta/internal/yaml.v2
- pkg/orderedmap
- pkg/filepos
pkg/yttlibraryext:
- pkg/yttlibraryext/toml
pkg/yttlibraryext/toml:
- pkg/yttlibrary
- pkg/orderedmap
- pkg/template/core
- pkg/yamlmeta
pkg/yttlibrary:
- pkg/version
- pkg/yttlibrary/overlay
- pkg/orderedmap
- pkg/template/core
- pkg/yamlmeta
pkg/yttlibrary/overlay:
- pkg/yamltemplate
- pkg/template
- pkg/filepos
- pkg/template/core
- pkg/yamlmeta
pkg/yamltemplate:
- pkg/texttemplate
- pkg/orderedmap
- pkg/template
- pkg/filepos
- pkg/template/core
- pkg/yamlmeta
pkg/texttemplate:
- pkg/template
- pkg/filepos
pkg/template:
- pkg/filepos
- pkg/template/core
pkg/template/core:
- pkg/orderedmap
pkg/cmd/template:
- pkg/workspace
- pkg/workspace/datavalues
- pkg/schema
- pkg/workspace/ref
- pkg/yttlibrary/overlay
- pkg/files
- pkg/cmd/ui
- pkg/template
- pkg/filepos
- pkg/yamlmeta
pkg/workspace:
- pkg/texttemplate
- pkg/workspace/datavalues
- pkg/yamltemplate
- pkg/yttlibrary
- pkg/schema
- pkg/workspace/ref
- pkg/yttlibrary/overlay
- pkg/files
- pkg/cmd/ui
- pkg/template
- pkg/filepos
- pkg/template/core
- pkg/yamlmeta
pkg/workspace/datavalues:
- pkg/schema
- pkg/workspace/ref
- pkg/template
- pkg/filepos
- pkg/template/core
- pkg/yamlmeta
pkg/schema:
- pkg/spell
- pkg/template
- pkg/filepos
- pkg/template/core
- pkg/yamlmeta
pkg/workspace/ref:
- pkg/template
- pkg/template/core
pkg/files:
- pkg/cmd/ui
main:
- pkg/cmd/template
- pkg/files
- pkg/cmd/ui
```
