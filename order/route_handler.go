package order

type orderRouteHandler struct {
	orderStore *orderStore
}

func NewRouteHandler() *orderRouteHandler {
	return &orderRouteHandler{
		orderStore: newOrderStore(),
	}
}
