# 📈 Sistema de Rebalanceo de Portafolio de Inversiones

![Go Version](https://img.shields.io/badge/Go-1.20%2B-00ADD8?style=flat&logo=go)
![Architecture](https://img.shields.io/badge/Architecture-Clean%20%2F%20Layered-blue)
![Tests](https://img.shields.io/badge/Tests-Passing-brightgreen)

Este proyecto implementa en **Go (Golang)** la lógica matemática y la orquestación para el **rebalanceo de un portafolio de inversión en acciones**. Permite calcular exactamente qué activos **comprar**, **vender** o **mantener** para alinear el valor real del portafolio del inversionista con su asignación de activos objetivo deseada.

---

## 🏛️ Arquitectura y Estructura del Proyecto

El código sigue los principios de **Clean Architecture** y separación clara de responsabilidades:

```text
fintual/
├── cmd/
│   └── main.go                  # Punto de entrada de la aplicación CLI (Entrypoint)
├── pkg/
│   ├── models/                  # Entidades y modelos del dominio
│   │   ├── stock.go             # Entidad Acción (Ticker, Cantidad, ObtenerPrecioActual)
│   │   ├── portfolio.go         # Entidad Portafolio y validación de asignación (100%)
│   │   ├── rebalance.go         # Estructuras de reporte, recomendaciones e iconos visuales
│   │   └── scenario.go          # Lector y mapeador de escenarios de prueba desde JSON
│   ├── services/                # Capa de servicio y lógica de negocio
│   │   ├── portfolio_service.go # Orquestador del algoritmo en 7 pasos (métodos privados)
│   │   ├── price_provider.go    # Interfaz y Mock para proveedores de precios de mercado
│   │   └── portfolio_service_test.go # Suite de pruebas unitarias automatizadas
│   └── cli/                     # Capa de presentación por consola (CLI)
│       ├── app.go               # Orquestador del menú interactivo y comandos CLI
│       └── printer.go           # Formateador e impresor visual de tablas y reportes
└── data/                        # Escenarios de prueba dinámicos en formato JSON
    ├── scenario1_tech_growth.json
    ├── scenario2_balanced.json
    └── scenario3_dividend_value.json
```

---

## ⚙️ Lógica del Algoritmo de Rebalanceo

La lógica central está encapsulada en `PortfolioService.Rebalance` y dividida en **7 pasos secuenciales**:

> [!IMPORTANT]
> **Regla de Validación**: La suma de los porcentajes de la asignación objetivo debe ser siempre exactamente igual al **100% (1.0)**.

| Paso | Nombre | Descripción / Fórmula |
| :--- | :--- | :--- |
| **Paso 1** | **Valor Total Actual** | Suma del valor de mercado de todas las acciones poseídas: <br> `Valor Total = ∑ (Cantidad × Precio Actual)` |
| **Paso 2** | **Asignación Objetivo** | Iteración sobre la distribución objetivo y validación del 100%. |
| **Paso 3** | **Valor Ideal** | Dinero ideal a tener invertido en cada acción: <br> `Valor Ideal = Valor Total Actual × Porcentaje Objetivo` |
| **Paso 4** | **Valor Real** | Dinero actual invertido en cada acción: <br> `Valor Real = Cantidad Poseída × Precio Actual` |
| **Paso 5** | **Determinar Diferencia** | Comparación entre lo ideal y lo real: <br> `Diferencia = Valor Ideal - Valor Real` <br> • `Diferencia > 0` ➔ **COMPRAR** <br> • `Diferencia < 0` ➔ **VENDER** <br> • `Diferencia == 0` ➔ **MANTENER** |
| **Paso 6** | **Traducir a Unidades** | Cálculo de acciones físicas a comprar o vender: <br> `Unidades = |Diferencia| / Precio Actual` |
| **Paso 7** | **Reporte Consolidado** | Generación del reporte completo ordenado con resumen de instrucciones. |

---

## 🚀 Instrucciones de Ejecución

> [!NOTE]
> Requisito previo: Tener instalado [Go 1.20+](https://go.dev/dl/).

### 1. Menú Interactivo (Recomendado)
Ejecuta la aplicación sin argumentos para abrir el menú interactivo de selección de escenarios:

```bash
go run cmd/main.go
```

**Vista en Consola:**
```text
================================================================================
             SISTEMA DE REBALANCEO DE PORTAFOLIO DE INVERSIONES                 
================================================================================

📋 Seleccione el escenario de prueba a ejecutar:

   [1] Portafolio Tech Growth
       └─ Rebalanceo de un portafolio de alto crecimiento tecnológico (scenario1_tech_growth.json)
   [2] Portafolio Balanceado Conservador
       └─ Distribución diversificada entre renta variable, bonos y oro (scenario2_balanced.json)
   [3] Portafolio Dividendos & Valor
       └─ Portafolio enfocado en generación de flujo de caja por dividendos (scenario3_dividend_value.json)
   [4] Procesar TODOS los escenarios secuencialmente

Ingrese el número de opción [Por defecto: 1]: 
```

> [!TIP]
> Si presionas `ENTER` sin ingresar ningún número, el sistema seleccionará automáticamente la opción **`[1]`** por defecto.

### 2. Ejecución Directa por CLI
Puedes enviar el número de opción directamente como argumento:

```bash
# Ejecutar escenario 2 (Balanceado Conservador)
go run cmd/main.go 2

# Ejecutar TODOS los escenarios dinámicamente
go run cmd/main.go 4
```

### 3. Cargar un Archivo JSON Específico
Puedes pasar la ruta directa de cualquier archivo JSON de portafolio:

```bash
go run cmd/main.go data/scenario3_dividend_value.json
```

---

## 🧪 Pruebas Unitarias

Para correr las pruebas unitarias automatizadas del paquete de servicios:

```bash
go test -v ./...
```

---

## 📄 Formato de Escenarios JSON (`data/*.json`)

Puedes agregar tus propios escenarios creando archivos en la carpeta `data/`:

```json
{
  "scenario_name": "Mi Escenario Personalizado",
  "description": "Prueba de rebalanceo personalizada",
  "market_prices": {
    "META": 500.00,
    "AAPL": 200.00,
    "GOOGL": 150.00
  },
  "portfolio": {
    "holdings": [
      { "ticker": "META", "quantity": 10.0 },
      { "ticker": "AAPL", "quantity": 50.0 }
    ],
    "target_allocation": {
      "META": 0.40,
      "AAPL": 0.40,
      "GOOGL": 0.20
    }
  }
}
```
