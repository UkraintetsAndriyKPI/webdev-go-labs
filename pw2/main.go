package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type EmissionInput struct {
	CoalAmount float64
	OilAmount  float64
	GasAmount  float64
}

type EmissionResult struct {
	CoalKTB float64
	CoalETB float64
	OilKTB  float64
	OilETB  float64
	GasKTB  float64
	GasETB  float64
}

type PageData struct {
	Input  EmissionInput
	Result string
}

var tmpl = template.Must(template.ParseFiles("./templates/index.html"))

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/calculate", calculate)

	fmt.Println("Server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	tmpl.Execute(w, PageData{})
}

func calculateEmission(input EmissionInput) EmissionResult {

	nCleaning := 0.985

	qrCoal := 20.47
	qrOil := 39.48
	qrGas := 33.08

	aOutCoal := 0.8
	aOutOil := 1.0
	aOutGas := 0.0

	AInCoal := 25.20
	AInOil := 0.15
	AInGas := 0.0

	GoutCoal := 1.5
	GoutOil := 0.0
	GoutGas := 0.0

	coalKTB := (math.Pow(10, 6) / qrCoal) *
		aOutCoal *
		(AInCoal / (100 - GoutCoal)) *
		(1 - nCleaning)

	coalETB := math.Pow(10, -6) *
		coalKTB *
		qrCoal *
		input.CoalAmount

	oilKTB := (math.Pow(10, 6) / qrOil) *
		aOutOil *
		(AInOil / (100 - GoutOil)) *
		(1 - nCleaning)

	oilETB := math.Pow(10, -6) *
		oilKTB *
		qrOil *
		input.OilAmount

	gasKTB := (math.Pow(10, 6) / qrGas) *
		aOutGas *
		(AInGas / (100 - GoutGas)) *
		(1 - nCleaning)

	gasETB := math.Pow(10, -6) *
		gasKTB *
		qrGas *
		input.GasAmount

	return EmissionResult{
		CoalKTB: coalKTB,
		CoalETB: coalETB,
		OilKTB:  oilKTB,
		OilETB:  oilETB,
		GasKTB:  gasKTB,
		GasETB:  gasETB,
	}
}

func calculate(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	parse := func(name string) float64 {
		v, _ := strconv.ParseFloat(r.FormValue(name), 64)
		return v
	}

	input := EmissionInput{
		CoalAmount: parse("coalAmount"),
		OilAmount:  parse("oilAmount"),
		GasAmount:  parse("gasAmount"),
	}

	resultStruct := calculateEmission(input)

	result := fmt.Sprintf(
		`Вугілля:
kₜᵦ = %.6f
Eₜᵦ = %.6f

Мазут:
kₜᵦ = %.6f
Eₜᵦ = %.6f

Газ:
kₜᵦ = %.6f
Eₜᵦ = %.6f`,
		resultStruct.CoalKTB,
		resultStruct.CoalETB,
		resultStruct.OilKTB,
		resultStruct.OilETB,
		resultStruct.GasKTB,
		resultStruct.GasETB,
	)

	data := PageData{
		Input:  input,
		Result: result,
	}

	tmpl.Execute(w, data)
}
