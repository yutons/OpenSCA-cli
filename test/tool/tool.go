package tool

import (
	"context"
	"testing"

	"github.com/xmirrorsecurity/opensca-cli/opensca"
	"github.com/xmirrorsecurity/opensca-cli/opensca/model"
	"github.com/xmirrorsecurity/opensca-cli/opensca/sca"
)

func Diff(a, b *model.DepGraph) bool {
	clear := func(p, n *model.DepGraph) bool {
		n.Path = ""
		n.Language = model.Lan_None
		return true
	}
	a.ForEachNode(clear)
	b.ForEachNode(clear)
	return a.Tree(false, true) != b.Tree(false, true)
}

func Dep(vendor, name, version string, children ...*model.DepGraph) *model.DepGraph {
	root := &model.DepGraph{
		Vendor:  vendor,
		Name:    name,
		Version: version,
	}
	for _, c := range children {
		root.AppendChild(c)
	}
	return root
}

func DevDep(vendor, name, version string, children ...*model.DepGraph) *model.DepGraph {
	root := Dep(vendor, name, version, children...)
	root.Develop = true
	return root
}

type TaskCase struct {
	Path   string
	Result *model.DepGraph
}

func RunTaskCase(t *testing.T, sca ...sca.Sca) func(cases []TaskCase) {
	return func(cases []TaskCase) {
		for _, c := range cases {
			deps, _ := opensca.RunTask(context.Background(), &opensca.TaskArg{
				DataOrigin: c.Path,
				Sca:        sca,
			})
			result := &model.DepGraph{}
			for _, dep := range deps {
				result.AppendChild(dep)
			}
			if Diff(result, c.Result) {
				t.Errorf("%s\nres:\n%sstd:\n%s", c.Path, result.Tree(false, true), c.Result.Tree(false, true))
			}
		}
	}
}
