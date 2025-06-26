package decorators

type RouteMapping struct {
	Method      string
	Path        string
	HandlerFunc any
}

type ControllerRegistration struct {
	BasePath string
	Instance any
}

var registeredControllers []ControllerRegistration
var registeredRoutes []RouteMapping

func RegisteredControllers() []ControllerRegistration {
	return registeredControllers
}

func RegisteredRoutes() []RouteMapping {
	return registeredRoutes
}

func Controller(basePath string, instance any) {
	registeredControllers = append(registeredControllers, ControllerRegistration{
		BasePath: basePath,
		Instance: instance,
	})
}

func GetMapping(path string, handler any) {
	registeredRoutes = append(registeredRoutes, RouteMapping{
		Method:      "GET",
		Path:        path,
		HandlerFunc: handler,
	})
}

func PostMapping(path string, handler any) {
	registeredRoutes = append(registeredRoutes, RouteMapping{
		Method:      "POST",
		Path:        path,
		HandlerFunc: handler,
	})
}

// Add PutMapping, DeleteMapping, etc. as needed
