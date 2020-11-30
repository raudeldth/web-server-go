package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "strconv"
)

var materias = make(map[string]map[string]float64)
var alumnos = make(map[string]map[string]float64)

func imprimeMaps() {
    fmt.Println("\n\tMATERIAS:\n", materias)
    fmt.Println("\n\n\tAlumnos:\n",alumnos)
}

func toString() (string, string) {
    var a string
    var m string

    for mat, alu := range materias {
        a += "<tr>" +
            "<th colspan='2'>" + mat + "</th>" +
            "</tr>" +
            "<tr><th>Alumno</th><th>Calificacion</th></tr>"
        for al, cal := range alu {
            a += "<tr>" +
                "<td>" + al + "</td>" +
                "<td>" + strconv.FormatFloat(cal, 'f', 2, 64) + "</td>" +
                "</tr>"
        }
    }
    for alu, ma := range alumnos {
        m += "<tr>" +
            "<th colspan='2'>" + alu + "</th>" +
            "</tr>" +
            "<tr><th>Materia</th><th>Calificacion</th></tr>"
        for mat, cal := range ma {
            m += "<tr>" +
                "<td>" + mat + "</td>" +
                "<td>" + strconv.FormatFloat(cal, 'f', 2, 64) + "</td>" +
                "</tr>"
        }
    }
    return a, m
}

func agregarAlumnoMateria(nombre string, materia string, cal float64 ) bool {
    alumno := make(map[string]float64)
    if _, ok := alumnos[nombre][materia]; ok {
        return false
    } else {
        alumno[nombre] = cal
        if _, ok := materias[materia] ; ok{
            materias[materia][nombre] = cal
        } else {
            materias[materia] = make(map[string]float64)
            materias[materia] = alumno
        }

        if _, ok := alumnos[nombre] ; ok{
            alumnos[nombre][materia] = cal
        } else {
            materiaM := make(map[string]float64)
            alumnos[nombre] = make(map[string]float64)
            materiaM[materia] = cal
            alumnos[nombre] = materiaM
        }
    }
    return true
}

func cargarHtml(a string) string {
    html, _ := ioutil.ReadFile(a)

    return string(html)
}

func form(res http.ResponseWriter, req *http.Request) {
    res.Header().Set(
        "Content-Type",
        "text/html",
    )
    fmt.Fprintf(
        res,
        cargarHtml("form.html"),
    )
}

func promedioGeneralAlumnos() float64{
    var cont float64
    var contAlu float64
    var totalAluCal float64
    var totalProms float64

    for _, element := range alumnos {
        for _, cal := range element {
            totalAluCal += cal
            cont++
        }
        totalProms += (totalAluCal/cont)
        totalAluCal = 0
        cont = 0
        contAlu++
    }
    return totalProms/contAlu
}

func promedioGeneral(res http.ResponseWriter, req *http.Request) {
    res.Header().Set(
        "Content-Type",
        "text/html",
    )

    fmt.Fprintf(
        res,
        cargarHtml("promedio-general.html"),
        strconv.FormatFloat(promedioGeneralAlumnos(), 'f', 2, 64),
    )
}

func selectMateria() string {
    var mat string
    for m, _ := range materias {
        mat += "<option value=" + m + ">" + m +
        "</option>"
    }

    return mat
}

func selectAlumno() string {
    var alu string
    for a, _ := range alumnos {
        alu += "<option value=" + a + ">" + a +
        "</option>"
    }

    return alu
}

func formMateria(res http.ResponseWriter, req *http.Request) {
    res.Header().Set(
        "Content-Type",
        "text/html",
    )

    fmt.Fprintf(
        res,
        cargarHtml("materia.html"),
        selectMateria(),
    )
}

func formAlumno(res http.ResponseWriter, req *http.Request) {
    res.Header().Set(
        "Content-Type",
        "text/html",
    )

    fmt.Fprintf(
        res,
        cargarHtml("alumno.html"),
        selectAlumno(),
    )
}

func promedioMateria(materia string) float64 {
    var cont float64
    var totalCal float64

    for key, element := range materias {
        if key == materia {
            for _, cal := range element {
                totalCal += cal
                cont++
            }
            break
        }
    }
    return totalCal/cont
}

func promedioAlumnos(alumno string) float64 {
    var cont float64
    var totalCal float64

    for key, element := range alumnos {
        if key == alumno {
            for _, cal := range element {
                totalCal += cal
                cont++
            }
            break
        }
    }
    return totalCal/cont
}

func materiaP(res http.ResponseWriter, req *http.Request) {
    fmt.Println(req.Method)
    switch req.Method {
    case "POST":
        if err := req.ParseForm(); err != nil {
            fmt.Fprintf(res, "ParseForm() error %v", err)
            return
        }

        res.Header().Set(
            "Content-Type",
            "text/html",
        )
        fmt.Fprintf(
            res,
            cargarHtml("promedio-materia.html"),
            strconv.FormatFloat(promedioMateria(req.FormValue("materia")), 'f', 2, 64),
        )
    case "GET":
        res.Header().Set(
            "Content-Type",
            "text/html",
        )
        a, _ := toString()
        fmt.Fprintf(
            res,
            cargarHtml("tabla-materia.html"),
            a,
        )
    }
}

func alumnoP(res http.ResponseWriter, req *http.Request) {
    fmt.Println(req.Method)
    switch req.Method {
    case "POST":
        if err := req.ParseForm(); err != nil {
            fmt.Fprintf(res, "ParseForm() error %v", err)
            return
        }

        res.Header().Set(
            "Content-Type",
            "text/html",
        )
        fmt.Fprintf(
            res,
            cargarHtml("promedio-alumno.html"),
            strconv.FormatFloat(promedioAlumnos(req.FormValue("alumno")), 'f', 2, 64),
        )
    case "GET":
        res.Header().Set(
            "Content-Type",
            "text/html",
        )
        _, m := toString()
        fmt.Fprintf(
            res,
            cargarHtml("tabla-alumno.html"),
            m,
        )
    }
}

func profesor(res http.ResponseWriter, req *http.Request) {
    fmt.Println(req.Method)
    switch req.Method {
    case "POST":
        if err := req.ParseForm(); err != nil {
            fmt.Fprintf(res, "ParseForm() error %v", err)
            return
        }
        fmt.Println(req.PostForm)
        cal, _ := strconv.ParseFloat(req.FormValue("calificacion"), 64)
        //fmt.Println(req.FormValue("calificacion"))
        reply := agregarAlumnoMateria(req.FormValue("alumno"), req.FormValue("materia"), cal)
        imprimeMaps()

        res.Header().Set(
            "Content-Type",
            "text/html",
        )

        al := req.FormValue("alumno")
        ma := req.FormValue("materia")
        if !reply {
            al += " no"
            ma += " no"
        }
        fmt.Fprintf(
            res,
            cargarHtml("respuesta.html"),
            al,
            ma,
        )
    case "GET":
        res.Header().Set(
            "Content-Type",
            "text/html",
        )
        a, m := toString()
        fmt.Fprintf(
            res,
            cargarHtml("tabla.html"),
            a,
            m,
        )
    }
}

func main() {
    http.HandleFunc("/form", form)
    http.HandleFunc("/profesor", profesor)
    http.HandleFunc("/materia", materiaP)
    http.HandleFunc("/alumno", alumnoP)
    http.HandleFunc("/formMateria", formMateria)
    http.HandleFunc("/formAlumno", formAlumno)
    http.HandleFunc("/general", promedioGeneral)
    fmt.Println("Corriendo servirdor de tareas...")
    http.ListenAndServe(":9000", nil)
}
