# 📈 Sistema de Rebalanceo de Portafolio de Inversiones

Este proyecto implementa en **Go (Golang)** la lógica matemática y de orquestación para el **rebalanceo de un portafolio de acciones**. El sistema permite determinar con precisión qué activos comprar, vender o mantener para alinear la cartera real del inversionista con su distribución objetivo deseada.

---

## 🏛️ Arquitectura y Estructura del Proyecto

El proyecto sigue los principios de **Clean Architecture** y separación de responsabilidades:

```text
fintual/
├── cmd/
│   └── main.go                  # Punto de entrada ultralimpio de la aplicación (14 líneas)
├── pkg/
│   ├── models/                  # Entidades del dominio
│   │   ├── stock.go             # Entidad Acción con método ObtenerPrecioActual
│   │   ├── portfolio.go         # Entidad Portafolio con validación de asignación (100%)
│   │   ├── rebalance.go         # Modelos de reporte, recomendaciones y ActionType.Symbol()
│   │   └── scenario.go          # Modelo de pruebas dinámicas y deserializador JSON
│   ├── services/                # Servicios y lógica de negocio
│   │   ├── portfolio_service.go # Orquestador desacoplado con métodos privados por cada paso
│   │   ├── price_provider.go    # Interfaz y Mock para consulta de precios de mercado
│   │   └── portfolio_service_test.go # Pruebas unitarias automatizadas
│   └── cli/                     # Componentes de interfaz por consola
│       ├── app.go               # Orquestador del menú interactivo y CLI
│       └── printer.go           # Formateador de tablas e instrucciones en consola
└── data/                        # Escenarios de prueba dinámicos en JSON
    ├── scenario1_tech_growth.json
    ├── scenario2_balanced.json
    └── scenario3_dividend_value.json
```

---

## ⚙️ Lógica del Algoritmo de Rebalanceo

La lógica central está encapsulada en `PortfolioService.Rebalance` y dividida en **7 pasos secuenciales**:

1. **Calcular el Valor Total Actual**: Suma el valor de mercado de todas las acciones poseídas ($\sum \text{cantidad} \times \text{precio\_actual}$).
2. **Iterar sobre la Asignación Objetivo**: Valida previamente que los porcentajes sumen exactamente **100% (1.0)**.
3. **Calcular el "Valor Ideal"**: Multiplica el Valor Total Actual por el porcentaje asignado a cada acción ($\text{Valor Total} \times \% \text{Objetivo}$).
4. **Calcular el "Valor Real"**: Consulta las unidades actualmente poseídas y las multiplica por su precio actual.
5. **Determinar la Diferencia**: Resta el Valor Real al Valor Ideal ($\text{Valor Ideal} - \text{Valor Real}$).
   - **Positivo**: Requiere **`COMPRAR`**.
   - **Negativo**: Requiere **`VENDER`**.
   - **Cero**: Se debe **`MANTENER`**.
6. **Traducir a Unidades (Acciones)**: Divide el monto de la diferencia entre el precio actual ($\frac{|\text{Diferencia}|}{\text{Precio Actual}}$).
7. **Retorno de Información**: Construye y devuelve el reporte consolidado con las acciones exactas a ejecutar.

---

## 🚀 Cómo Ejecutar la Aplicación

### Prerrequisitos
- Tener instalado **Go** (versión 1.20 o superior).

### 1. Menú Interactivo (Modo por Defecto)
Ejecuta el comando principal sin argumentos para abrir el menú interactivo:

```bash
go run cmd/main.go
```

**Ejemplo de consola:**
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

*(Si presionas `ENTER` sin escribir nada, seleccionará automáticamente la opción `[1]`)*.

### 2. Ejecución Directa por Número de Opción
Puedes pasar el número de opción como argumento en la terminal:

```bash
# Ejecutar escenario 2 (Balanceado Conservador)
go run cmd/main.go 2

# Ejecutar TODOS los escenarios
go run cmd/main.go 4
```

### 3. Ejecutar un Archivo JSON Específico
Puedes pasar la ruta directa a cualquier archivo JSON:

```bash
go run cmd/main.go data/scenario3_dividend_value.json
```

---

## 🧪 Cómo Ejecutar las Pruebas Unitarias

Para correr el suite completo de pruebas unitarias automatizadas:

```bash
go test -v ./...
```

---

## 📄 Estructura de Escenarios JSON (`data/*.json`)

Puedes agregar nuevos escenarios creando archivos `.json` en la carpeta `data/` siguiendo esta estructura:

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
