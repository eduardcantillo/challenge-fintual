package cli

import (
	"fmt"

	"fintual/pkg/models"
)

// PrintReport presenta en consola el reporte formateado con tabla e instrucciones finales.
func PrintReport(report *models.RebalanceReport) {
	fmt.Printf("💰 VALOR TOTAL ACTUAL DEL PORTAFOLIO: $%.2f USD\n", report.TotalPortfolioValue)
	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("%-8s | %-10s | %-12s | %-12s | %-12s | %-12s | %-10s\n",
		"TICKER", "ACCIÓN", "CANT. ACTUAL", "VALOR REAL", "VALOR IDEAL", "DIFERENCIA", "UNIDADES")
	fmt.Println("--------------------------------------------------------------------------------")

	for _, item := range report.Items {
		fmt.Printf("%-8s | %-10s | %12.2f | $%11.2f | $%11.2f | $%11.2f | %10.4f\n",
			item.Ticker,
			item.Action.Symbol(),
			item.CurrentQuantity,
			item.RealValue,
			item.IdealValue,
			item.DifferenceAmount,
			item.UnitsToTrade,
		)
	}

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Println()
	fmt.Println("📌 RESUMEN DE INSTRUCCIONES:")
	for _, item := range report.Items {
		switch item.Action {
		case models.ActionBuy, models.ActionSell:
			fmt.Printf("   • %s %.4f acciones de %s (Precio actual: $%.2f)\n",
				item.Action, item.UnitsToTrade, item.Ticker, item.CurrentPrice)
		default:
			fmt.Printf("   • %s posición en %s (Totalmente balanceada)\n",
				item.Action, item.Ticker)
		}
	}
	fmt.Println("================================================================================")
}
