package registry

import (
	"microtecture/infrastructure/application"
	"microtecture/infrastructure/datastore"
	"microtecture/interface/controllers"
	uc "microtecture/usecase/controllers"
)

// Registry interface
type Registry interface {
	NewRootController() uc.Root
}

type registry struct {
	controller     application.Controller
	restController application.RestController
}

// New creates and returns registry
func New() (Registry, error) {
	app, err := application.New()
	if err != nil {
		return nil, err
	}

	ctrl, err := application.NewController(app)
	if err != nil {
		return nil, err
	}
	restController := application.NewRestController(ctrl)

	return registry{ctrl, restController}, nil
}

// NewTestRegistry creates and return registry for test goals
func NewTest() (Registry, error) {
	app, err := application.New()
	if err != nil {
		return nil, err
	}

	session, err := datastore.NewTestSession()
	if err != nil {
		return nil, err
	}
	app.DBSession = *session

	c, err := application.NewController(app)
	if err != nil {
		return nil, err
	}
	rc := application.NewRestController(c)

	return registry{c, rc}, nil
}

// NewRootController creates and return root controller
func (self registry) NewRootController() uc.Root {
	apiv1 := controllers.NewApiV1(self.restController)

	root := controllers.NewRoot(self.restController, apiv1)

	return root
}
