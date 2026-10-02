package services_test

import (
	"math"
	"testing"

	"fintual/pkg/models"
	"fintual/pkg/services"
)

func TestRebalance_Success(t *testing.T) {
	// Precios simulados
	prices := map[string]float64{
		"META": 500.0,
		"AAPL": 200.0,
	}
	priceProvider := services.NewMockPriceProvider(prices)

	// Portafolio actual:
	// 10 acciones de META @ $500 = $5,000
	// 50 acciones de AAPL @ $200 = $10,000
	// Valor Total Actual = $15,000
	portfolio := models.NewPortfolio()
	portfolio.AddStock(models.NewStock("META", 10))
	portfolio.AddStock(models.NewStock("AAPL", 50))

	// Distribución deseada: 40% META, 60% AAPL
	portfolio.SetTargetAllocation("META", 0.40)
	portfolio.SetTargetAllocation("AAPL", 0.60)

	portfolioService := services.NewPortfolioService(priceProvider)
	report, err := portfolioService.Rebalance(portfolio)

	if err != nil {
		t.Fatalf("se esperaba un rebalanceo exitoso, pero ocurrió error: %v", err)
	}

	// Verificar Valor Total: $15,000
	if math.Abs(report.TotalPortfolioValue-15000.0) > 1e-4 {
		t.Errorf("Valor Total devuelto insatisfactorio: se esperaba 15000, se obtuvo %.2f", report.TotalPortfolioValue)
	}

	// META:
	// Valor Ideal: $15,000 * 0.40 = $6,000
	// Valor Real: 10 * $500 = $5,000
	// Diferencia: +$1,000 -> COMPRAR $1,000 / $500 = 2 unidades
	var metaItem models.RebalanceItem
	var aaplItem models.RebalanceItem
	for _, item := range report.Items {
		if item.Ticker == "META" {
			metaItem = item
		} else if item.Ticker == "AAPL" {
			aaplItem = item
		}
	}

	if metaItem.Action != models.ActionBuy {
		t.Errorf("Para META se esperaba COMPRAR, se obtuvo %s", metaItem.Action)
	}
	if math.Abs(metaItem.UnitsToTrade-2.0) > 1e-4 {
		t.Errorf("Para META se esperaban 2 unidades a comprar, se obtuvo %.4f", metaItem.UnitsToTrade)
	}

	// AAPL:
	// Valor Ideal: $15,000 * 0.60 = $9,000
	// Valor Real: 50 * $200 = $10,000
	// Diferencia: -$1,000 -> VENDER $1,000 / $200 = 5 unidades
	if aaplItem.Action != models.ActionSell {
		t.Errorf("Para AAPL se esperaba VENDER, se obtuvo %s", aaplItem.Action)
	}
	if math.Abs(aaplItem.UnitsToTrade-5.0) > 1e-4 {
		t.Errorf("Para AAPL se esperaban 5 unidades a vender, se obtuvo %.4f", aaplItem.UnitsToTrade)
	}
}

func TestRebalance_InvalidAllocation(t *testing.T) {
	priceProvider := services.NewMockPriceProvider(map[string]float64{"META": 500.0})
	portfolio := models.NewPortfolio()
	portfolio.AddStock(models.NewStock("META", 10))

	// Porcentajes que no suman 100% (suman 90%)
	portfolio.SetTargetAllocation("META", 0.90)

	portfolioService := services.NewPortfolioService(priceProvider)
	_, err := portfolioService.Rebalance(portfolio)

	if err == nil {
		t.Error("se esperaba un error por asignación inválida != 100%, pero no ocurrió ninguno")
	}
}
