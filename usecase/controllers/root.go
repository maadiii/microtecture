package controllers

import "microtecture/infrastructure/application"

// Root is root controller interface
type Root interface {
	GetBase() application.RestController
}

// ApiV1 is api v1 controller interface
type ApiV1 interface{}
