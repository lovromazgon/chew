package funcmap

import (
	"errors"
	"fmt"
	"sync"
)

var (
	Global Functions
	initOnce  sync.Once
)

func initFunctions() {
	initOnce.Do(func() {
		Global = make([]*Func, 0)
	})
}

func AddFunc(f *Func) {
	initFunctions()
	Global = Global.MustAddFunc(f)
}

type Func struct {
	Func interface{}
	Doc  FuncDoc
}

type FuncDoc struct {
	Name     string
	Text     string
	Example  string
	Template string

	NestedFuncs []FuncDoc
}

func NewFunc(fun interface{}, name, text, example string, nestedFuncs ...FuncDoc) *Func {
	f := &Func{
		Func: fun,
		Doc: FuncDoc{
			Name:        name,
			Text:        text,
			Example:     example,
			NestedFuncs: nestedFuncs,
		},
	}

	if len(nestedFuncs) == 0 {
		f.Doc.Template = "simpleFuncDoc"
	} else {
		f.Doc.Template = "nestedFuncDoc"
	}

	return f
}

func NewNestedFuncDoc(name, text, example string) FuncDoc {
	return FuncDoc{
		Name:        name,
		Text:        text,
		Example:     example,
		Template:    "simpleFuncDoc",
		NestedFuncs: nil,
	}
}

type Functions []*Func

func (fs Functions) FuncMap() map[string]interface{} {
	funcMap := make(map[string]interface{})

	for _, f := range fs {
		funcMap[f.Doc.Name] = f.Func
	}

	return funcMap
}

func (fs Functions) MustAddFunc(f *Func) Functions {
	fsNew, err := fs.AddFunc(f)
	if err != nil {
		panic(err)
	}
	return fsNew
}

func (fs Functions) AddFunc(f *Func) (Functions, error) {
	for _, fExisting := range fs {
		if f.Doc.Name == fExisting.Doc.Name {
			return nil, errors.New(fmt.Sprintf("Function map already contains function %s", f.Doc.Name))
		}
	}
	return append(fs, f), nil
}
