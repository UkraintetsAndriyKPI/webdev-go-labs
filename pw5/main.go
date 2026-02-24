package main

import (
	"fmt"
	"html/template"
	"log"
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

	electricGasSwitchAmount := parseFloat(r, "electricGasSwitchAmount")
	electricPL100kBLen := parseFloat(r, "electricPL100kBLen")
	converter110kBAmount := parseFloat(r, "converter110kBAmount")
	switcher10kBAmount := parseFloat(r, "switcher10kBAmount")
	connector10kBAmount := parseFloat(r, "connector10kBAmount")

	electricGasSwitchFRate := 0.01
	electricPL100kBFRate := 0.007
	converter110kBFRate := 0.015
	switcher10kBFRate := 0.02
	connector10kBFRate := 0.03

	failureRate1wheel :=
		electricGasSwitchFRate*electricGasSwitchAmount +
			electricPL100kBFRate*electricPL100kBLen +
			converter110kBFRate*converter110kBAmount +
			switcher10kBFRate*switcher10kBAmount +
			connector10kBFRate*connector10kBAmount

	var avgRepairTime float64
	if failureRate1wheel != 0 {
		avgRepairTime = failureRate1wheel / failureRate1wheel
	}

	accidentCoeff := failureRate1wheel * avgRepairTime / 8760
	idleCoeff := 1.2 * 43 / 8760

	failureRate2wheel :=
		2*failureRate1wheel*(accidentCoeff+idleCoeff) +
			switcher10kBFRate

	conclusion := "Одноколова система електропередачі надійніша"
	if failureRate1wheel > failureRate2wheel {
		conclusion = "Двоколова система електропередачі надійніша"
	}

	result := fmt.Sprintf(
		`Вхідні дані:

Елегазові вимикачі: %.2f
ПЛ-100 кВ (довжина): %.2f
Трансформатори 110 кВ: %.2f
Вимикачі 10 кВ: %.2f
З'єднувачі 10 кВ: %.2f

----------------------------------------

Результати:

Частота відмови одноколової системи: %.6f рік^-1
Частота відмови двоколової системи: %.6f рік^-1

Висновок:
%s`,
		electricGasSwitchAmount,
		electricPL100kBLen,
		converter110kBAmount,
		switcher10kBAmount,
		connector10kBAmount,
		failureRate1wheel,
		failureRate2wheel,
		conclusion,
	)

	tmpl.Execute(w, PageData{
		Result:    result,
		ActiveTab: "calc1",
	})
}


func calculate2Handler(w http.ResponseWriter, r *http.Request) {

	log.Println("POST /calculate2")
	r.ParseForm()

	accidentCost := parseFloat(r, "accidentCost")
	scheduleCost := parseFloat(r, "scheduleCost")
	denyRate := parseFloat(r, "denyRate")
	avgAccidentDenyTime := parseFloat(r, "avgAccidentDenyTime")
	avgScheduleDenyTime := parseFloat(r, "avgScheduleDenyTime")

	mathHopeAccident :=
		denyRate * avgAccidentDenyTime * 5120 * 6451

	mathHopeSchedule :=
		avgScheduleDenyTime * 5120 * 6451

	mathHopeAccidentCost :=
		mathHopeAccident * accidentCost

	mathHopeScheduleCost :=
		mathHopeSchedule * scheduleCost

	mathHopeAtAllCost :=
		mathHopeAccidentCost + mathHopeScheduleCost

	result := fmt.Sprintf(
		`Вхідні дані:

Вартість аварійного переривання: %.2f грн/кВт·год
Вартість планового переривання: %.2f грн/кВт·год
Частота відмов: %.6f
Середній час аварійного відключення: %.2f год
Середній час планового відключення: %.2f год

----------------------------------------

Результати:

Математичне сподівання аварійного невідпущення: %.2f кВт*год
Математичне сподівання планового невідпущення: %.2f кВт*год

Збитки від аварійного переривання: %.2f грн
Збитки від планового переривання: %.2f грн

Загальні збитки: %.2f грн`,
		accidentCost,
		scheduleCost,
		denyRate,
		avgAccidentDenyTime,
		avgScheduleDenyTime,
		mathHopeAccident,
		mathHopeSchedule,
		mathHopeAccidentCost,
		mathHopeScheduleCost,
		mathHopeAtAllCost,
	)

	tmpl.Execute(w, PageData{
		Result:    result,
		ActiveTab: "calc2",
	})
}
