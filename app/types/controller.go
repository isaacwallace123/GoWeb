package types

// Controller interface
type Controller interface {
	BasePath() string
	Routes() []Route
}

// ControllerBase: minimal struct for declarative controller setup
type ControllerBase struct {
	basePath string
	routes   []Route
}

// Fluent builder for BasePath
func (c *ControllerBase) WithBasePath(path string) *ControllerBase {
	c.basePath = path
	return c
}

// Fluent builder for Routes
func (c *ControllerBase) WithRoutes(routes []Route) *ControllerBase {
	c.routes = routes
	return c
}

// Implements Controller interface
func (c *ControllerBase) BasePath() string { return c.basePath }
func (c *ControllerBase) Routes() []Route  { return c.routes }
