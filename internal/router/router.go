package router

type Router struct {
	routes map[string]interface{}
}

func New() *Router {
	return &Router{
		routes: make(map[string]interface{}),
	}
}
