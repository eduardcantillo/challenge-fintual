package models

// PriceProvider defines the contract for fetching the latest stock prices.
type PriceProvider interface {
	GetPrice(ticker string) (float64, error)
}

// Stock representa cada activo individual dentro del portafolio.
type Stock struct {
	Ticker   string  `json:"ticker"`   // Nombre o identificador de la acción (ej. META, AAPL)
	Quantity float64 `json:"quantity"` // Cantidad de unidades (acciones) actualmente poseídas
}

// NewStock crea una nueva instancia de Stock.
func NewStock(ticker string, quantity float64) Stock {
	return Stock{
		Ticker:   ticker,
		Quantity: quantity,
	}
}

// ObtenerPrecioActual es el método que devuelve el último precio disponible
// de esa acción específica utilizando el proveedor de precios.
func (s Stock) ObtenerPrecioActual(provider PriceProvider) (float64, error) {
	return provider.GetPrice(s.Ticker)
}
