package govisual

import (
	"errors"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"
)

func New(c Config) (*Visualization, error) {
	if c.Module == "" {
		c.Module = guessModuleName(c.RootPath)
	}
	if c.Module == "" {
		return nil, errors.New("cannot guess the module name")
	}
	return &Visualization{Config: c, Digraph: NewDigraph()}, nil
}

type ImportType uint

const (
	Sys ImportType = iota
	Third
	Org
	Self
)

type Package struct {
	FullName     string
	SimpleName   string
	Type         ImportType
	RelativePath string
}

func (p Package) Hash() string {
	return p.FullName
}

type Config struct {
	RootPath   string // project root path
	TargetPath string // target directory path, relative to RootPath
	Module     string // project module name
	All        bool   // parse all subdirectories
	SysPkg     bool   // include go sdk packages
	ThirdPkg   bool   // include third party packages
	OrgPkg     bool   // include same organization packages
	HostAsOrg  bool   // use host as organization
}

type Visualization struct {
	Config  Config
	Digraph *Digraph
}

func (p *Visualization) Parse(root string) error {
	dir, err := ioutil.ReadDir(root)
	if err != nil {
		return err
	}

	if p.Config.All {
		for _, fi := range dir {
			if strings.HasPrefix(fi.Name(), ".") || !fi.IsDir() || fi.Name() == "vendor" {
				continue
			}

			if err := p.Parse(pathjoin(root, fi.Name())); err != nil {
				return err
			}
		}
	}

	return p.parse(root)
}

func (p *Visualization) parse(rpath string) error {
	pkgs, err := parser.ParseDir(token.NewFileSet(), pathjoin(p.Config.RootPath, rpath),
		sourceFilter, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		if rpath == "." {
			rpath = pkg.Name
		}
		from := &Package{
			FullName:   pathjoin(p.Config.Module, rpath),
			SimpleName: rpath,
			Type:       Self,
		}
		for _, file := range pkg.Files {
			for _, imp := range file.Imports {

				to := p.decodeImport(imp.Path.Value)

				if (to.Type == Sys && !p.Config.SysPkg) ||
					(to.Type == Third && !p.Config.ThirdPkg) ||
					(to.Type == Org && !p.Config.OrgPkg) {
					continue
				}

				p.Digraph.Edge(from, to)
				if p.Digraph.HasFrom(to.Hash()) {
					continue
				}

				if to.Type == Self && to.RelativePath != "" {
					if err := p.parse(to.RelativePath); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (p *Visualization) decodeImport(s string) *Package {
	s = unquote(s)
	pkg := &Package{
		FullName:   s,
		SimpleName: s,
		Type:       Sys,
	}

	parts1 := strings.Split(pkg.FullName, "/")
	host1 := parts1[0]

	defer func() {
		if pkg.Type == Self {
			if p.Config.HostAsOrg {
				pkg.RelativePath = pathjoin(parts1[2:]...)
			} else {
				pkg.RelativePath = pathjoin(parts1[3:]...)
			}
			pkg.SimpleName = pkg.RelativePath
		}

		if p.Config.OrgPkg && (pkg.Type == Self || pkg.Type == Org) {
			if p.Config.HostAsOrg {
				pkg.SimpleName = pathjoin(parts1[1:]...)
			} else {
				pkg.SimpleName = pathjoin(parts1[2:]...)
			}
		}

	}()

	// host not contains '.', it's a sys pkg
	if !strings.Contains(host1, ".") {
		return pkg
	}

	pkg.Type = Third
	parts2 := strings.Split(p.Config.Module, "/")
	host2 := parts2[0]

	// hosts not equal, it's a third party pkg
	if host1 != host2 {
		return pkg
	}

	orgIndex := 1
	if p.Config.HostAsOrg {
		orgIndex = 0
	}

	// hosts equals, but organization name not equal, it's a third party pkg
	if parts1[orgIndex] != parts2[orgIndex] {
		return pkg
	}

	if p.Config.OrgPkg {
		pkg.Type = Org
	}

	repo1 := parts1[orgIndex+1]
	repo2 := parts2[orgIndex+1]
	//  organization name equals, but repository name not equal, it's a organization pkg or a third party pkg
	if repo1 != repo2 {
		return pkg
	}

	// repository name equals, it's a pkg in this repo
	pkg.Type = Self
	return pkg
}
