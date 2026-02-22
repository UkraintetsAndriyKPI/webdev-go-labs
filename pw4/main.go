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
	Result    string
	ActiveTab string
	Form      map[string]float64
}

var tmpl = template.Must(template.ParseFiles("./templates/index.html"))

func main() {

	log.Println("====================================")
	log.Println("Server is starting on http://localhost:8080")
	log.Println("====================================")

	http.HandleFunc("/", home)
	http.HandleFunc("/calculate1", calculate1Handler)
	http.HandleFunc("/calculate2", calculate2Handler)
	http.HandleFunc("/calculate3", calculate3Handler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	log.Println("GET / - Home page opened")
	tmpl.Execute(w, PageData{ActiveTab: "calc1"})
}

func parseFloat(r *http.Request, name string) float64 {
	value := r.FormValue(name)
	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Printf("Error parsing field %s (value: %s)\n", name, value)
		return 0
	}
	return v
}


func calculate1Handler(w http.ResponseWriter, r *http.Request) {

	log.Println("POST /calculate1")
	r.ParseForm()

	electricCurrent := parseFloat(r, "electricCurrent")
	load := parseFloat(r, "load")

	log.Printf("Input -> ElectricCurrent: %.4f | Load: %.4f\n",
		electricCurrent, load)

	voltage := 10.0

	normalModeCurrent := load / (math.Sqrt(3) * 2 * voltage)
	emergencyModeCurrent := 2 * normalModeCurrent
	sek := normalModeCurrent / 1.4
	smin := electricCurrent * 1000

	log.Printf("Result -> Normal: %.4f | Emergency: %.4f | Sek: %.4f | Smin: %.4f\n",
		normalModeCurrent, emergencyModeCurrent, sek, smin)

	result := fmt.Sprintf(
		`Вхідні дані:
Electric Current = %.4f
Load = %.4f

Результати:
Тип: ААБ кабель
Струм норм. режим: %.4f
Струм авар. режим: %.4f
Економічний переріз: %.4f
Мінімальний переріз: %.4f`,
		electricCurrent,
		load,
		normalModeCurrent,
		emergencyModeCurrent,
		sek,
		smin,
	)

	tmpl.Execute(w, PageData{
		Result:    result,
		ActiveTab: "calc1",
		Form: map[string]float64{
			"electricCurrent": electricCurrent,
			"load":            load,
		},
	})
}


func calculate2Handler(w http.ResponseWriter, r *http.Request) {

	log.Println("POST /calculate2")
	r.ParseForm()

	powerMBA := parseFloat(r, "powerMBA")
	log.Printf("Input -> PowerMBA: %.4f\n", powerMBA)

	voltage := 10.5

	Xc := math.Pow(voltage, 2) / powerMBA
	Xt := voltage/100 * math.Pow(voltage, 2) / 6.3
	Xsum := Xc + Xt
	Ipo := voltage / (math.Sqrt(3) * Xsum)

	log.Printf("Result -> Xc: %.4f | Xt: %.4f | Xsum: %.4f | Ipo: %.4f\n",
		Xc, Xt, Xsum, Ipo)

	result := fmt.Sprintf(
		`Вхідні дані:
Power (MVA) = %.4f

Результати:
Xc = %.4f
Xt = %.4f
Xsum = %.4f
Iп0 = %.4f`,
		powerMBA,
		Xc, Xt, Xsum, Ipo,
	)

	tmpl.Execute(w, PageData{
		Result:    result,
		ActiveTab: "calc2",
		Form: map[string]float64{
			"powerMBA": powerMBA,
		},
	})
}

func calculate3Handler(w http.ResponseWriter, r *http.Request) {

	log.Println("POST /calculate3")
	r.ParseForm()

	resistanceMax := parseFloat(r, "resistanceMax")
	resistanceBH := parseFloat(r, "resistanceBH")
	valueRcn := parseFloat(r, "valueRcn")
	valueXcn := parseFloat(r, "valueXcn")
	valueRcnMIN := parseFloat(r, "valueRcnMIN")
	valueXcnMIN := parseFloat(r, "valueXcnMIN")

	log.Printf("Input -> Rmax: %.4f | Voltage: %.4f | Rcn: %.4f | Xcn: %.4f | RcnMin: %.4f | XcnMin: %.4f\n",
		resistanceMax,
		resistanceBH,
		valueRcn,
		valueXcn,
		valueRcnMIN,
		valueXcnMIN)

	reactiveResistance :=
		(resistanceMax * math.Pow(resistanceBH, 2)) /
			(100 * 6.3)

	busX := valueXcn + reactiveResistance
	busZ := math.Sqrt(math.Pow(valueRcn, 2) + math.Pow(busX, 2))

	busXmin := valueXcnMIN + reactiveResistance
	busZmin := math.Sqrt(math.Pow(valueRcnMIN, 2) + math.Pow(busXmin, 2))

	stream3 := (resistanceBH * 1000) / (math.Sqrt(3) * busZ)
	stream2 := stream3 * math.Sqrt(3) / 2
	stream3min := (resistanceBH * 1000) / (math.Sqrt(3) * busZmin)
	stream2min := stream3min * math.Sqrt(3) / 2

	log.Printf("Result -> Xt: %.4f | I3: %.4f | I2: %.4f | I3min: %.4f | I2min: %.4f\n",
		reactiveResistance,
		stream3,
		stream2,
		stream3min,
		stream2min)

	result := fmt.Sprintf(
		`Вхідні дані:
ResistanceMax = %.4f
Voltage = %.4f
Rcn = %.4f
Xcn = %.4f
RcnMin = %.4f
XcnMin = %.4f

Результати:
Xт = %.4f
Iш(3) норм = %.4f
Iш(2) норм = %.4f
Iш(3) мін = %.4f
Iш(2) мін = %.4f`,
		resistanceMax,
		resistanceBH,
		valueRcn,
		valueXcn,
		valueRcnMIN,
		valueXcnMIN,
		reactiveResistance,
		stream3,
		stream2,
		stream3min,
		stream2min,
	)

	tmpl.Execute(w, PageData{
		Result:    result,
		ActiveTab: "calc3",
	})
}
