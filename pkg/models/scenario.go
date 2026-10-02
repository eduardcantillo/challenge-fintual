package models

import (
	"encoding/json"
	"fmt"
	"os"
)

// Scenario representa un caso de prueba o simulación dinámica cargado desde un archivo JSON.
type Scenario struct {
	Name          string             `json:"scenario_name"`
	Description   string             `json:"description"`
	MarketPrices  map[string]float64 `json:"market_prices"`
	PortfolioData struct {
		Holdings         []Stock            `json:"holdings"`
		TargetAllocation map[string]float64 `json:"target_allocation"`
	} `json:"portfolio"`
}

// LoadScenarioFromFile lee y deserializa un archivo JSON a una estructura Scenario.
func LoadScenarioFromFile(filePath string) (*Scenario, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error al leer el archivo JSON '%s': %w", filePath, err)
	}

	var scenario Scenario
	if err := json.Unmarshal(bytes, &scenario); err != nil {
		return nil, fmt.Errorf("error al interpretar JSON de '%s': %w", filePath, err)
	}

	return &scenario, nil
}

// ToPortfolio convierte los datos deserializados en una entidad Portfolio de nuestro dominio.
func (s *Scenario) ToPortfolio() *Portfolio {
	portfolio := NewPortfolio()

	for _, stock := range s.PortfolioData.Holdings {
		portfolio.AddStock(stock)
	}

	for ticker, percentage := range s.PortfolioData.TargetAllocation {
		portfolio.SetTargetAllocation(ticker, percentage)
	}

	return portfolio
}
