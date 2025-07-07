package types

type Controller interface {
	BasePath() string
	Routes() []Route
}

type ControllerBase struct {
	basePath       string
	routes         []Route
	preMiddleware  []MiddlewareFunc
	postMiddleware []MiddlewareFunc
}

// WithBasePath will set the URI of the controller (Like "/api/v1/users")
func (c *ControllerBase) WithBasePath(path string) *ControllerBase {
	c.basePath = path
	return c
}

// WithRoutes adds the routes that will be handled by the controller created
func (c *ControllerBase) WithRoutes(routes []Route) *ControllerBase {
	c.routes = routes
	return c
}

// Use adds pre-middleware (runs before handler)
func (c *ControllerBase) Use(mw ...MiddlewareFunc) *ControllerBase {
	c.preMiddleware = append(c.preMiddleware, mw...)
	return c
}

// UseAfter adds post-middleware (runs after handler, before global post-middleware)
func (c *ControllerBase) UseAfter(mw ...MiddlewareFunc) *ControllerBase {
	c.postMiddleware = append(c.postMiddleware, mw...)
	return c
}

// BasePath, Routes, PreMiddleware, & PostMiddleware These are the pre-determined methods that users are essentially FORCED to use because GoWeb uses these methods in the dispatch
func (c *ControllerBase) BasePath() string                 { return c.basePath }
func (c *ControllerBase) Routes() []Route                  { return c.routes }
func (c *ControllerBase) PreMiddleware() []MiddlewareFunc  { return c.preMiddleware }
func (c *ControllerBase) PostMiddleware() []MiddlewareFunc { return c.postMiddleware }
