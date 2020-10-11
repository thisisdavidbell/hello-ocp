package controller

import (
	"github.com/thisisdavidbell/hello-ocp/hello-ocp-operator/pkg/controller/helloocp"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, helloocp.Add)
}
