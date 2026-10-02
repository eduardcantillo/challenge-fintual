package models

// ActionType representa el tipo de acción a realizar durante el rebalanceo.
type ActionType string

const (
	ActionBuy  ActionType = "COMPRAR"
	ActionSell ActionType = "VENDER"
	ActionHold ActionType = "MANTENER"
)

// Symbol devuelve la etiqueta visual formateada con su icono representativo.
func (a ActionType) Symbol() string {
	switch a {
	case ActionBuy:
		return "🟢 COMPRAR"
	case ActionSell:
		return "🔴 VENDER"
	default:
		return "⚪ MANTENER"
	}
}

// RebalanceItem representa el resultado del cálculo para una acción individual.
type RebalanceItem struct {
	Ticker           string     `json:"ticker"`
	Action           ActionType `json:"action"`
	UnitsToTrade     float64    `json:"units_to_trade"`
	CurrentPrice     float64    `json:"current_price"`
	CurrentQuantity  float64    `json:"current_quantity"`
	RealValue        float64    `json:"real_value"`
	IdealValue       float64    `json:"ideal_value"`
	DifferenceAmount float64    `json:"difference_amount"`
	TargetPercentage float64    `json:"target_percentage"`
}

// RebalanceReport representa el reporte consolidado del portafolio rebalanceado.
type RebalanceReport struct {
	TotalPortfolioValue float64         `json:"total_portfolio_value"`
	Items               []RebalanceItem `json:"items"`
}
