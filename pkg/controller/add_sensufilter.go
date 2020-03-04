package controller

import (
	"github.com/betorvs/sensu-operator/pkg/controller/sensufilter"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, sensufilter.Add)
}
