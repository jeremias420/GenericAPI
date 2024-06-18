package app

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

var Plantillas *template.Template

// METODOS
func Inicio(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Generic API")
}

func Index(w http.ResponseWriter, r *http.Request) {
	Plantillas.ExecuteTemplate(w, "index", nil)
}

func ObtenerSucesos(w http.ResponseWriter, r *http.Request) {
	// Ejecutar la consulta SQL
	query := "SELECT * FROM suceso"
	ObtenerConsultaJSON(w, r, query)
}

func ObtenerLocalidades(w http.ResponseWriter, r *http.Request) {
	// Ejecutar la consulta SQL
	query := "SELECT * FROM localidad"
	ObtenerConsultaJSON(w, r, query)
}

func ObtenerSuceso(w http.ResponseWriter, r *http.Request) {
	//obtener parametros
	vars := mux.Vars(r)
	suceID := vars["suceID"]

	// Ejecutar la consulta SQL
	if suceID != "" {
		query := "SELECT * FROM suceso where suce_id = " + suceID
		ObtenerConsultaJSON(w, r, query)
	}
}

func InsertarSala(w http.ResponseWriter, r *http.Request) {
	// Lee el cuerpo de la solicitud HTTP en un []byte.
	body := GetBody(r)

	if body != nil {
		descripcion := ObtenerParametroPOST(body, "descripcion")
		areaID := ObtenerParametroPOST(body, "areaID")

		if descripcion != "" && areaID != "" {
			query := "exec sp_insertarSala " + Comillas(descripcion) + " , " + areaID
			ObtenerConsultaJSON(w, r, query)
		} else {
			http.Error(w, "Usuario o Password invalidos!", http.StatusBadRequest)
		}

	} else {
		http.Error(w, "No se puede leer el body", http.StatusBadRequest)
	}
}

func myfunc() {

}
