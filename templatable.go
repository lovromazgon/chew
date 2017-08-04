package chew

type Templatable interface {
	// returns the filename of the template to use (without .tmpl)
	Template() string
	// returns object which implements this interface (self)
	Data() interface{}
}

// GoTemplatable is the default implementation of Templatable which can be
// extended and used in other objects
type GoTemplatable struct {
	template string
	data     interface{}
}

func (gt GoTemplatable) Template() string {
	return gt.template
}

func (gt GoTemplatable) Data() interface{} {
	return gt.data
}