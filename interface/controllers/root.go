package controllers

import (
	"microtecture/infrastructure/application"
	"microtecture/usecase/controllers"
)

type root struct {
	application.RestController
	ApiV1 controllers.ApiV1
}

// NewRootController creates and returns root controller
func NewRoot(c application.RestController, apiv1 controllers.ApiV1) controllers.Root {
	return root{c, apiv1}
}

func (self root) GetBase() application.RestController {
	return self.RestController
}

type apiv1 struct {
	application.RestController
}

// NewApiv1Controller creates and returns apiv1 controller
func NewApiV1(c application.RestController) controllers.ApiV1 {
	return apiv1{c}
}
