package models

import (
	"fmt"
	"math"
)

// Portfolio representa la clase contenedora de las acciones poseídas
// y de la asignación objetivo (porcentajes deseados).
type Portfolio struct {
	Holdings         map[string]Stock   `json:"holdings"`          // Colección de Acciones poseídas (Ticker -> Stock)
	TargetAllocation map[string]float64 `json:"target_allocation"` // Asignación Objetivo (Ticker -> Porcentaje en decimal, ej: 0.40 = 40%)
}

// NewPortfolio crea e inicializa un nuevo Portafolio.
func NewPortfolio() *Portfolio {
	return &Portfolio{
		Holdings:         make(map[string]Stock),
		TargetAllocation: make(map[string]float64),
	}
}

// AddStock agrega o actualiza una acción dentro del portafolio.
func (p *Portfolio) AddStock(stock Stock) {
	p.Holdings[stock.Ticker] = stock
}

// SetTargetAllocation establece la asignación objetivo deseada para un ticker.
func (p *Portfolio) SetTargetAllocation(ticker string, percentage float64) {
	p.TargetAllocation[ticker] = percentage
}

// ValidateAllocation valida que la suma de todos los porcentajes de asignación
// sea exactamente 100% (1.0). Devuelve error si la validación falla.
func (p *Portfolio) ValidateAllocation() error {
	if len(p.TargetAllocation) == 0 {
		return fmt.Errorf("la asignación objetivo no puede estar vacía")
	}

	var totalPercentage float64
	for ticker, percentage := range p.TargetAllocation {
		if percentage < 0 {
			return fmt.Errorf("el porcentaje para la acción %s no puede ser negativo: %.2f%%", ticker, percentage*100)
		}
		totalPercentage += percentage
	}

	// Tolerancia épsilon para imprecisiones de coma flotante
	const epsilon = 1e-4
	if math.Abs(totalPercentage-1.0) > epsilon {
		return fmt.Errorf("validación fallida: la suma de las asignaciones objetivo debe ser exactamente 100%% (1.0). Suma actual: %.2f%%", totalPercentage*100)
	}

	return nil
}
