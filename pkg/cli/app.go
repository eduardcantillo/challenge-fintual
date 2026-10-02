package cli

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"fintual/pkg/models"
	"fintual/pkg/services"
)

// MenuOption representa cada item seleccionable del menú de escenarios.
type MenuOption struct {
	FilePath string
	Scenario *models.Scenario
}

// CLIApp orquesta la interfaz de usuario por consola y la ejecución de escenarios.
type CLIApp struct {
	dataDir string
}

// NewCLIApp crea una nueva instancia de la aplicación CLI.
func NewCLIApp(dataDir string) *CLIApp {
	return &CLIApp{dataDir: dataDir}
}

// Run inicia el flujo principal de la aplicación CLI.
func (app *CLIApp) Run(args []string) error {
	app.printHeader()

	options, err := app.loadOptions()
	if err != nil {
		return err
	}

	if len(args) > 1 {
		return app.handleCLIArg(args[1], options)
	}

	return app.runInteractiveMenu(options)
}

func (app *CLIApp) printHeader() {
	fmt.Println("================================================================================")
	fmt.Println("             SISTEMA DE REBALANCEO DE PORTAFOLIO DE INVERSIONES                 ")
	fmt.Println("================================================================================")
	fmt.Println()
}

func (app *CLIApp) loadOptions() ([]MenuOption, error) {
	pattern := filepath.Join(app.dataDir, "*.json")
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		return nil, fmt.Errorf("no se encontraron archivos JSON en '%s'", pattern)
	}

	var options []MenuOption
	for _, path := range matches {
		sc, err := models.LoadScenarioFromFile(path)
		if err == nil {
			options = append(options, MenuOption{FilePath: path, Scenario: sc})
		}
	}

	if len(options) == 0 {
		return nil, fmt.Errorf("no se pudo cargar ningún escenario desde '%s'", app.dataDir)
	}

	return options, nil
}

func (app *CLIApp) handleCLIArg(arg string, options []MenuOption) error {
	if idx, err := strconv.Atoi(arg); err == nil {
		if idx >= 1 && idx <= len(options) {
			return app.ExecuteOption(options[idx-1])
		}
		if idx == len(options)+1 {
			return app.ExecuteAllOptions(options)
		}
	}
	return app.ExecuteFromFile(arg)
}

func (app *CLIApp) runInteractiveMenu(options []MenuOption) error {
	fmt.Println("📋 Seleccione el escenario de prueba a ejecutar:")
	fmt.Println()

	for i, opt := range options {
		fmt.Printf("   [%d] %s\n", i+1, opt.Scenario.Name)
		fmt.Printf("       └─ %s (%s)\n", opt.Scenario.Description, filepath.Base(opt.FilePath))
	}
	fmt.Printf("   [%d] Procesar TODOS los escenarios secuencialmente\n", len(options)+1)
	fmt.Println()

	const defaultChoice = 1
	fmt.Printf("Ingrese el número de opción [Por defecto: %d]: ", defaultChoice)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	choice := defaultChoice
	if input != "" {
		if val, err := strconv.Atoi(input); err == nil {
			choice = val
		} else {
			fmt.Printf("\n⚠️  Entrada no válida ('%s'). Usando opción por defecto [%d].\n", input, defaultChoice)
		}
	}

	fmt.Println()

	if choice >= 1 && choice <= len(options) {
		return app.ExecuteOption(options[choice-1])
	} else if choice == len(options)+1 {
		return app.ExecuteAllOptions(options)
	}

	fmt.Printf("⚠️  Opción fuera de rango (%d). Ejecutando opción por defecto [%d]...\n\n", choice, defaultChoice)
	return app.ExecuteOption(options[defaultChoice-1])
}

func (app *CLIApp) ExecuteOption(opt MenuOption) error {
	fmt.Printf("📂 Cargando escenario desde '%s'...\n", opt.FilePath)
	return app.ExecuteScenario(opt.Scenario)
}

func (app *CLIApp) ExecuteAllOptions(options []MenuOption) error {
	for idx, opt := range options {
		fmt.Printf("📂 ESCENARIO #%d: Cargando desde '%s'...\n", idx+1, opt.FilePath)
		if err := app.ExecuteScenario(opt.Scenario); err != nil {
			return err
		}
		fmt.Println()
	}
	return nil
}

func (app *CLIApp) ExecuteFromFile(filePath string) error {
	sc, err := models.LoadScenarioFromFile(filePath)
	if err != nil {
		return err
	}
	return app.ExecuteScenario(sc)
}

func (app *CLIApp) ExecuteScenario(scenario *models.Scenario) error {
	fmt.Printf("🎯 Escenario: %s\n", scenario.Name)
	fmt.Printf("📝 Descripción: %s\n", scenario.Description)
	fmt.Println()

	priceProvider := services.NewMockPriceProvider(scenario.MarketPrices)
	portfolio := scenario.ToPortfolio()

	portfolioService := services.NewPortfolioService(priceProvider)
	report, err := portfolioService.Rebalance(portfolio)
	if err != nil {
		return fmt.Errorf("error al ejecutar rebalanceo: %w", err)
	}

	PrintReport(report)
	return nil
}
