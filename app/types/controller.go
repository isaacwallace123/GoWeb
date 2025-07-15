package types

type Controller interface {
	BasePath() string
	Routes() []Route
}

type ControllerBase struct {
	basePath       string
	routes         []Route
	preMiddleware  []Middleware
	postMiddleware []Middleware
}

// WithBasePath sets the base path for the controller (e.g. "/api/users")
func (c *ControllerBase) WithBasePath(path string) *ControllerBase {
	c.basePath = path
	return c
}

// WithRoutes adds the route list for this controller
func (c *ControllerBase) WithRoutes(routes []Route) *ControllerBase {
	c.routes = routes
	return c
}

// Use adds pre-handler middleware
func (c *ControllerBase) Use(mw ...Middleware) *ControllerBase {
	c.preMiddleware = append(c.preMiddleware, mw...)
	return c
}

// UseAfter adds post-handler middleware
func (c *ControllerBase) UseAfter(mw ...Middleware) *ControllerBase {
	c.postMiddleware = append(c.postMiddleware, mw...)
	return c
}

// Required interface implementations
func (c *ControllerBase) BasePath() string             { return c.basePath }
func (c *ControllerBase) Routes() []Route              { return c.routes }
func (c *ControllerBase) PreMiddleware() []Middleware  { return c.preMiddleware }
func (c *ControllerBase) PostMiddleware() []Middleware { return c.postMiddleware }
