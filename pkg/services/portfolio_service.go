package services

import (
	"fmt"
	"math"
	"sort"

	"fintual/pkg/models"
)

// PortfolioService es el servicio orquestador que contiene la lógica de negocio del portafolio.
type PortfolioService struct {
	priceProvider models.PriceProvider
}

// NewPortfolioService crea una nueva instancia del servicio de portafolio.
func NewPortfolioService(priceProvider models.PriceProvider) *PortfolioService {
	return &PortfolioService{
		priceProvider: priceProvider,
	}
}

// Rebalance es el método principal que orquesta los pasos del rebalanceo mediante métodos privados.
func (s *PortfolioService) Rebalance(portfolio *models.Portfolio) (*models.RebalanceReport, error) {
	// Validar previamente que el portafolio y los porcentajes objetivos sean válidos (Suma = 100%)
	if err := portfolio.ValidateAllocation(); err != nil {
		return nil, fmt.Errorf("error de validación de asignación objetivo: %w", err)
	}

	// Paso 1: Calcular el Valor Total Actual del portafolio
	totalPortfolioValue, currentPrices, err := s.calcularValorTotalActual(portfolio)
	if err != nil {
		return nil, err
	}

	// Recopilar tickers únicos ordenados
	tickers := s.obtenerTickersAEvaluar(portfolio)

	// Pasos 2 al 7: Iterar sobre la asignación objetivo y calcular acciones requeridas
	items, err := s.procesarRebalanceoAcciones(portfolio, tickers, totalPortfolioValue, currentPrices)
	if err != nil {
		return nil, err
	}

	return &models.RebalanceReport{
		TotalPortfolioValue: totalPortfolioValue,
		Items:               items,
	}, nil
}

// -----------------------------------------------------------------------------
// Métodos Internos (Privados) - Dividen el servicio en responsabilidades pequeñas
// -----------------------------------------------------------------------------

// Paso 1: Recorre la colección de acciones y calcula el valor total actual en dinero.
func (s *PortfolioService) calcularValorTotalActual(portfolio *models.Portfolio) (float64, map[string]float64, error) {
	var totalValue float64
	currentPrices := make(map[string]float64)

	for ticker, stock := range portfolio.Holdings {
		price, err := stock.ObtenerPrecioActual(s.priceProvider)
		if err != nil {
			return 0, nil, fmt.Errorf("error al obtener precio actual para '%s': %w", ticker, err)
		}
		currentPrices[ticker] = price
		totalValue += stock.Quantity * price
	}

	if totalValue <= 0 {
		return 0, nil, fmt.Errorf("el valor total del portafolio es $0.00. Se requieren activos con valor para rebalancear")
	}

	return totalValue, currentPrices, nil
}

// Paso 2 (Auxiliar): Obtiene la lista unificada y ordenada de tickers a procesar.
func (s *PortfolioService) obtenerTickersAEvaluar(portfolio *models.Portfolio) []string {
	tickerSet := make(map[string]bool)
	for ticker := range portfolio.TargetAllocation {
		tickerSet[ticker] = true
	}
	for ticker := range portfolio.Holdings {
		tickerSet[ticker] = true
	}

	var tickers []string
	for ticker := range tickerSet {
		tickers = append(tickers, ticker)
	}
	sort.Strings(tickers)
	return tickers
}

// Obtener o consultar precio actual de una acción específica utilizando ObtenerPrecioActual
func (s *PortfolioService) obtenerPrecioActualAccion(ticker string, currentPrices map[string]float64) (float64, error) {
	if price, exists := currentPrices[ticker]; exists {
		return price, nil
	}

	dummyStock := models.NewStock(ticker, 0)
	price, err := dummyStock.ObtenerPrecioActual(s.priceProvider)
	if err != nil {
		return 0, fmt.Errorf("error al obtener precio actual para '%s': %w", ticker, err)
	}
	currentPrices[ticker] = price
	return price, nil
}

// Paso 3: Multiplica el Valor Total Actual por el porcentaje asignado a la acción.
func (s *PortfolioService) calcularValorIdeal(totalPortfolioValue, targetPercentage float64) float64 {
	return totalPortfolioValue * targetPercentage
}

// Paso 4: Consulta la cantidad actual poseída y la multiplica por su precio actual.
func (s *PortfolioService) calcularValorReal(stockHolding models.Stock, holdingExists bool, price float64) (float64, float64) {
	var currentQuantity float64
	if holdingExists {
		currentQuantity = stockHolding.Quantity
	}
	realValue := currentQuantity * price
	return currentQuantity, realValue
}

// Paso 5: Resta el Valor Real al Valor Ideal para determinar diferencia y si se debe COMPRAR, VENDER o MANTENER.
func (s *PortfolioService) determinarDiferenciaYAccion(idealValue, realValue float64) (float64, models.ActionType) {
	differenceAmount := idealValue - realValue
	const epsilonAmount = 1e-4

	if differenceAmount > epsilonAmount {
		return differenceAmount, models.ActionBuy
	} else if differenceAmount < -epsilonAmount {
		return differenceAmount, models.ActionSell
	}
	return 0, models.ActionHold
}

// Paso 6: Toma la diferencia calculada y la divide por el precio actual para obtener unidades a operar.
func (s *PortfolioService) traducirAUnidades(differenceAmount, price float64, action models.ActionType) float64 {
	if action == models.ActionHold || price <= 0 {
		return 0
	}
	return math.Abs(differenceAmount) / price
}

// Pasos 2 al 7: Orquesta el procesamiento de cada ticker individual utilizando los métodos privados anteriores.
func (s *PortfolioService) procesarRebalanceoAcciones(
	portfolio *models.Portfolio,
	tickers []string,
	totalPortfolioValue float64,
	currentPrices map[string]float64,
) ([]models.RebalanceItem, error) {
	var items []models.RebalanceItem

	for _, ticker := range tickers {
		targetPercentage := portfolio.TargetAllocation[ticker]

		price, err := s.obtenerPrecioActualAccion(ticker, currentPrices)
		if err != nil {
			return nil, err
		}

		// Paso 3: Calcular el Valor Ideal
		idealValue := s.calcularValorIdeal(totalPortfolioValue, targetPercentage)

		// Paso 4: Calcular el Valor Real
		stockHolding, holdingExists := portfolio.Holdings[ticker]
		currentQuantity, realValue := s.calcularValorReal(stockHolding, holdingExists, price)

		// Paso 5: Determinar la Diferencia y la Acción
		differenceAmount, action := s.determinarDiferenciaYAccion(idealValue, realValue)

		// Paso 6: Traducir a Unidades
		unitsToTrade := s.traducirAUnidades(differenceAmount, price, action)

		// Paso 7: Guardar el resultado en el reporte
		item := models.RebalanceItem{
			Ticker:           ticker,
			Action:           action,
			UnitsToTrade:     unitsToTrade,
			CurrentPrice:     price,
			CurrentQuantity:  currentQuantity,
			RealValue:        realValue,
			IdealValue:       idealValue,
			DifferenceAmount: differenceAmount,
			TargetPercentage: targetPercentage,
		}

		items = append(items, item)
	}

	return items, nil
}
