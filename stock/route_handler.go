package stock

type stockRouteHandler struct {
	stockStore *stockStore
}

func NewRouteHandler() *stockRouteHandler {
	return &stockRouteHandler{
		stockStore: newStockStore(),
	}
}
