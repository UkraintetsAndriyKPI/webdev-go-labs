package main

import (
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"strconv"
)

type PageData struct {
	Result string
	Table  [][]string
	Form   map[string]string
}

var tmpl = template.Must(template.ParseFiles("./templates/index.html"))

var KKDnominal = 0.92
var powerCoefficient = 0.9
var loadVoltage = 0.38

var amountEP = map[string]float64{
	"GrindingMachine":    4,
	"DrillingMachine":    2,
	"GroutingMachine":    4,
	"CircularSaw":        1,
	"Press":              1,
	"PolishingMachine":   1,
	"MillingMachine":     2,
	"Fan":                1,
	"WeldingTransformer": 2,
	"DryerWardrobe":      2,
}

var nominalPower = map[string]float64{
	"GrindingMachine":    20,
	"DrillingMachine":    14,
	"GroutingMachine":    42,
	"CircularSaw":        36,
	"Press":              20,
	"PolishingMachine":   40,
	"MillingMachine":     32,
	"Fan":                20,
	"WeldingTransformer": 100,
	"DryerWardrobe":      120,
}

var usageCoeff = map[string]float64{
	"GrindingMachine":    0.15,
	"DrillingMachine":    0.12,
	"GroutingMachine":    0.15,
	"CircularSaw":        0.3,
	"Press":              0.5,
	"PolishingMachine":   0.2,
	"MillingMachine":     0.2,
	"Fan":                0.65,
	"WeldingTransformer": 0.2,
	"DryerWardrobe":      0.8,
}

var reactivePowerCoeff = map[string]float64{
	"GrindingMachine":    1.33,
	"DrillingMachine":    1.0,
	"GroutingMachine":    1.33,
	"CircularSaw":        1.52,
	"Press":              0.75,
	"PolishingMachine":   1.0,
	"MillingMachine":     1.0,
	"Fan":                0.75,
	"WeldingTransformer": 3.0,
	"DryerWardrobe":      0.0,
}

func main() {
	log.Println("====================================")
	log.Println("Server is starting on http://localhost:8080")
	log.Println("====================================")

	http.HandleFunc("/", home)
	http.HandleFunc("/calculate", calculateHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, PageData{
		Form: make(map[string]string),
	})
}

func parseFloat(value string) (float64, bool) {
	if value == "" {
		return 0, false
	}
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, false
	}
	return v, true
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	formValues := make(map[string]string)

	for key, val := range r.Form {
		if len(val) > 0 {
			formValues[key] = val[0]
		}
	}

	// оновлення параметрів якщо введені
	if v, ok := parseFloat(r.FormValue("nominalPowerGrindingMachine")); ok {
		nominalPower["GrindingMachine"] = v
	}
	if v, ok := parseFloat(r.FormValue("usageCoeffPolishingMachine")); ok {
		usageCoeff["PolishingMachine"] = v
	}
	if v, ok := parseFloat(r.FormValue("reactivePowerCoeffCircularSaw")); ok {
		reactivePowerCoeff["CircularSaw"] = v
	}

	stream1stLevel := make(map[string]float64)
	nPnKB := make(map[string]float64)
	nPnKBtg := make(map[string]float64)
	nPnpow2 := make(map[string]float64)

	for key := range amountEP {
		stream1stLevel[key] = amountEP[key] * nominalPower[key]
		nPnKB[key] = amountEP[key] * nominalPower[key] * usageCoeff[key]
		nPnKBtg[key] = nPnKB[key] * reactivePowerCoeff[key]
		nPnpow2[key] = amountEP[key] * math.Pow(nominalPower[key], 2)
	}

	var sumStream, sumPnKB, sumPnKBtg, sumPnpow2 float64
	for key := range stream1stLevel {
		sumStream += stream1stLevel[key]
		sumPnKB += nPnKB[key]
		sumPnKBtg += nPnKBtg[key]
		sumPnpow2 += nPnpow2[key]
	}

	effectiveEPamount := math.Pow(sumStream-stream1stLevel["WeldingTransformer"]-stream1stLevel["DryerWardrobe"], 2) / sumPnpow2
	calculatedActiveLoad := 1.25 * sumPnKB
	fullPower := math.Sqrt(math.Pow(calculatedActiveLoad, 2) + math.Pow(sumPnKBtg, 2))

	// додаткові обчислення для шин та цеху
	cehEPEfectiveAmount := sumPnKB / sumStream
	cehUsageCoeff := sumPnKB / sumStream
	cehUsageCoeffAtAll := 1.25 * cehUsageCoeff
	calculatedActiveLoadOnTires := cehUsageCoeffAtAll * sumStream
	calculatedReactiveLoadOnTires := sumPnKBtg
	fullPowerOnTires := math.Sqrt(math.Pow(calculatedActiveLoadOnTires, 2) + math.Pow(calculatedReactiveLoadOnTires, 2))
	calculatedCurrentOnTires := fullPowerOnTires / loadVoltage / math.Sqrt(3)

	result := fmt.Sprintf(
		`Ефективна к-сть ЕП ШР1=ШР2=ШР3: %.2f
Розрахунковий коеф. актив. потуж. ШР1=ШР2=ШР3: 1.25
Розрахункове актив. навантаження ШР1=ШР2=ШР3: %.2f
Розрахункове реактивне навантаження ШР1=ШР2=ШР3: %.2f
Повна потужність ШР1=ШР2=ШР3: %.2f
----------------------------------------
Ефективна к-сть ЕП: %.2f
Коеф. використання всього цеху: %.2f
Розрахунковий коеф. актив. потужності цеху вцілому: %.2f
Розрахункове актив. навантаження на шинах 0.38 кВ ТП: %.2f кВт
Розрахункове реактив. навантаження на шинах 0.38 кВ ТП: %.2f квар
Повна потужність на шинах 0.38 кВ ТП: %.2f кВА
Розрахунковий груп. струм на шинах 0.38 кВ ТП: %.2f А`,
		effectiveEPamount,
		calculatedActiveLoad,
		sumPnKBtg,
		fullPower,
		cehEPEfectiveAmount,
		cehUsageCoeff,
		cehUsageCoeffAtAll,
		calculatedActiveLoadOnTires,
		calculatedReactiveLoadOnTires,
		fullPowerOnTires,
		calculatedCurrentOnTires,
	)

	var table [][]string
	for key := range amountEP {
		row := []string{
			key,
			fmt.Sprintf("%.0f", amountEP[key]),
			fmt.Sprintf("%.2f", stream1stLevel[key]),
			fmt.Sprintf("%.2f", usageCoeff[key]),
			fmt.Sprintf("%.2f", nPnKB[key]),
			fmt.Sprintf("%.2f", nPnKBtg[key]),
			fmt.Sprintf("%.2f", nPnpow2[key]),
		}
		table = append(table, row)
	}

	tmpl.Execute(w, PageData{
		Result: result,
		Table:  table,
		Form:   formValues,
	})
}
