package types

type RouteEntry struct {
	Method  string
	Path    string
	Handler string
}

type Controller interface {
	BasePath() string
	Routes() []RouteEntry
}
