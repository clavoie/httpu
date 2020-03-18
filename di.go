package httpu

import "github.com/clavoie/di/v2"

// NewDiDefs returns a new collection of di definitions that
// can be used to inject httpu into your project.
func NewDiDefs() []*di.Def {
	return []*di.Def{
		{NewImpl, di.PerHttpRequest},
	}
}
