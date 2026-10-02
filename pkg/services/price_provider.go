package services

import (
	"fmt"
	"sync"
)

// MockPriceProvider es una implementación en memoria de PriceProvider para pruebas y simulación.
type MockPriceProvider struct {
	mu     sync.RWMutex
	prices map[string]float64
}

// NewMockPriceProvider crea una nueva instancia con precios iniciales opcionales.
func NewMockPriceProvider(initialPrices map[string]float64) *MockPriceProvider {
	provider := &MockPriceProvider{
		prices: make(map[string]float64),
	}
	for ticker, price := range initialPrices {
		provider.prices[ticker] = price
	}
	return provider
}

// SetPrice actualiza o asigna el precio de una acción.
func (m *MockPriceProvider) SetPrice(ticker string, price float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.prices[ticker] = price
}

// GetPrice obtiene el último precio disponible para una acción específica.
func (m *MockPriceProvider) GetPrice(ticker string) (float64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	price, exists := m.prices[ticker]
	if !exists {
		return 0, fmt.Errorf("no se encontró precio para la acción '%s'", ticker)
	}
	if price <= 0 {
		return 0, fmt.Errorf("precio inválido (%.2f) para la acción '%s'", price, ticker)
	}
	return price, nil
}
