package payment

type paymentRouteHandler struct {
	paymentStore *paymentStore
}

func NewRouteHandler() *paymentRouteHandler {
	return &paymentRouteHandler{
		paymentStore: newPaymentStore(),
	}
}
