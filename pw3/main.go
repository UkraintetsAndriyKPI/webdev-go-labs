package main

import (
	"fmt"
	"html/template"
	"math"
	"net/http"
	"strconv"
)

type FormInput struct {
	AvgDayPower       float64
	StandardDeviation float64
	ElectricityPrice  float64
}

type PageData struct {
	Input  FormInput
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

func calculate(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	parse := func(name string) float64 {
		v, _ := strconv.ParseFloat(r.FormValue(name), 64)
		return v
	}

	input := FormInput{
		AvgDayPower:       parse("avgDayPower"),
		StandardDeviation: parse("standardDeviation"),
		ElectricityPrice:  parse("electricityPrice"),
	}

	result := calculateLogic(input)

	data := PageData{
		Input:  input,
		Result: result,
	}

	tmpl.Execute(w, data)
}

func calculateLogic(data FormInput) string {

	// Before improvements
	omegaW1 := integration(4.75, 5.25, 8, data) / 100
	electChangeW11 := 24 * data.AvgDayPower * omegaW1
	electChangeW12 := 24 * data.AvgDayPower * (1 - omegaW1)
	profitP1 := 24 * data.AvgDayPower * omegaW1 * data.ElectricityPrice
	fine1 := 24 * data.AvgDayPower * (1 - omegaW1) * data.ElectricityPrice
	clearProfitP1 := profitP1 - fine1

	var conclusion1 string
	if clearProfitP1 > 0 {
		conclusion1 = "Electro station is profitable. Positive profit"
	} else {
		conclusion1 = "Electro station is not profitable. Negative profit"
	}

	// After improvements
	omegaW2 := integration(4.75, 5.25, 26, data) / 100
	electChangeW21 := 24 * data.AvgDayPower * omegaW2
	electChangeW22 := 24 * data.AvgDayPower * (1 - omegaW2)
	profitP2 := 24 * data.AvgDayPower * omegaW2 * data.ElectricityPrice
	fine2 := 24 * data.AvgDayPower * (1 - omegaW2) * data.ElectricityPrice
	clearProfitP2 := profitP2 - fine2

	var conclusion2 string
	if clearProfitP2 > 0 {
		conclusion2 = "Electro station is profitable. Positive profit"
	} else {
		conclusion2 = "Electro station is not profitable. Negative profit"
	}

	return fmt.Sprintf(
		`Before improvements:
δW1 = %.6f
W1 = %.6f
W2 = %.6f
profit = %.6f
fine = %.6f
clear profit = %.6f
conclusion = %s

After improvements:
δW2 = %.6f
W1 = %.6f
W2 = %.6f
profit = %.6f
fine = %.6f
clear profit = %.6f
conclusion = %s`,
		omegaW1,
		electChangeW11,
		electChangeW12,
		profitP1,
		fine1,
		clearProfitP1,
		conclusion1,
		omegaW2,
		electChangeW21,
		electChangeW22,
		profitP2,
		fine2,
		clearProfitP2,
		conclusion2,
	)
}

func integration(start, end float64, steps int, data FormInput) float64 {
	sum := 0.0
	step := (end - start) / float64(steps)

	i := start
	for i < end {
		sum += expression(i, data)
		i += step
	}

	return sum
}

func expression(p float64, data FormInput) float64 {
	return (1 / data.StandardDeviation * math.Sqrt(2*math.Pi)) *
		math.Pow(math.E,
			math.Pow(p-data.AvgDayPower, 2) /
				2 *
				math.Pow(data.StandardDeviation, 2),
		)
}
