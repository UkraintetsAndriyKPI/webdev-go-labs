package main

import (
    "fmt"
    "html/template"
    "net/http"
    "strconv"
)

type PageData struct {
    Tab    string
    Result string
}

var tmpl = template.Must(template.ParseFiles("./templates/index.html"))

func main() {
    http.HandleFunc("/", home)
    http.HandleFunc("/calculate", calculate)
    http.HandleFunc("/calculate2", calculate2)

    http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
    tab := r.URL.Query().Get("tab")
    if tab != "1" && tab != "2" {
        tab = "1"
    }

    data := PageData{
        Tab: tab,
    }

    tmpl.Execute(w, data)
}
func calculate(w http.ResponseWriter, r *http.Request) {

    r.ParseForm()

    parse := func(name string) float64 {
        v, _ := strconv.ParseFloat(r.FormValue(name), 64)
        return v
    }

    HP := parse("HP")
    CP := parse("CP")
    SP := parse("SP")
    NP := parse("NP")
    OP := parse("OP")
    WP := parse("WP")
    AP := parse("AP")

    // Коефіцієнти
    KPC := 100 / (100 - WP)
    KPG := 100 / (100 - WP - AP)

    // Перераховані показники (C)
    HC := HP * KPC
    CC := CP * KPC
    SC := SP * KPC
    NC := NP * KPC
    OC := OP * KPC
    AC := AP * KPC

    // Перераховані показники (G)
    HG := HP * KPG
    CG := CP * KPG
    SG := SP * KPG
    NG := NP * KPG
    OG := OP * KPG

    // Формули теплотворності
    QPH := (339*CP +
        1030*HP -
        108.8*(OP-SP) -
        25*WP) / 1000

    QCH := (QPH + 0.025*WP) * KPC
    QGH := (QPH + 0.025*WP) * KPG

    result := fmt.Sprintf(
        `KPC = %.4f
KPG = %.4f

HC = %.4f
CC = %.4f
SC = %.4f
NC = %.4f
OC = %.4f
AC = %.4f

HG = %.4f
CG = %.4f
SG = %.4f
NG = %.4f
OG = %.4f

QPH = %.4f
QCH = %.4f
QGH = %.4f`,
        KPC, KPG,
        HC, CC, SC, NC, OC, AC,
        HG, CG, SG, NG, OG,
        QPH, QCH, QGH,
    )

    data := PageData{
        Tab:    "1",
        Result: result,
    }

    tmpl.Execute(w, data)
}

func calculate2(w http.ResponseWriter, r *http.Request) {

    r.ParseForm()

    parse := func(name string) float64 {
        v, _ := strconv.ParseFloat(r.FormValue(name), 64)
        return v
    }

    HG := parse("HG")
    CG := parse("CG")
    SG := parse("SG")
    OG := parse("OG")
    WG := parse("WG")
    AG := parse("AG")
    QI := parse("QI")

    // Перевірка на коректність
    if WG+AG >= 100 {
        http.Error(w, "Invalid input values", http.StatusBadRequest)
        return
    }

    coefficient := (100 - WG - AG) / 100

    // Обчислення
    HP := HG * coefficient
    CP := CG * coefficient
    SP := SG * coefficient
    OP := OG * coefficient
    AP := (AG * (100 - WG)) / 100

    QPI := QI*coefficient - 0.025*WG

    result := fmt.Sprintf(
        `HP = %.4f
CP = %.4f
SP = %.4f
OP = %.4f
AP = %.4f

QPI = %.4f`,
        HP, CP, SP, OP, AP, QPI,
    )

    data := PageData{
        Tab:    "2",
        Result: result,
    }

    tmpl.Execute(w, data)
}
