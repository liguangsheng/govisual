package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/liguangsheng/govisual"
)

type StringArray []string

func (i *StringArray) String() string {
	return strings.Join(*i, "|")
}

func (i *StringArray) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func main() {
	var templates StringArray
	var hostAsOrg string
	config := govisual.Config{}
	flag.BoolVar(&config.SysPkg, "sys", false, "include go sdk packages")
	flag.BoolVar(&config.ThirdPkg, "third", false, "include third party packages")
	flag.BoolVar(&config.OrgPkg, "org", false, "include all subdirectory")
	flag.BoolVar(&config.HostAsOrg, "host-as-org", false, "include all subdirectory")
	flag.BoolVar(&config.All, "all", false, "include all subdirectory")
	flag.StringVar(&config.Module, "module", "", "custom module name")
	flag.StringVar(&hostAsOrg, "hostasorg", "", "include organization packages")
	flag.Var(&templates, "templates", "custom template files")
	flag.Parse()

	target := flag.Arg(0)
	if target == "" {
		target = "."
	}
	config.TargetPath = target

	absRootPath, err := filepath.Abs(".")
	check(err)
	config.RootPath = absRootPath

	v, err := govisual.New(config)
	check(err)
	check(v.Parse(target))

	fmt.Println(govisual.Render(v, templates...))
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
