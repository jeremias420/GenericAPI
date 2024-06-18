package main

import (
	"GenericAPI/app"
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	Config "GenericAPI/app/config"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("Starting application...")

	app.Plantillas = template.Must(template.ParseGlob("templates/*"))
	//cargar configuracion
	err := Config.LoadConfig()
	if err != nil {
		logrus.Error("main.LoadConfig.error: " + err.Error())
		return
	}
	router := mux.NewRouter()

	fmt.Println("Servidor corriendo...")

	//definicion de rutas
	router.HandleFunc("/", app.Inicio)                                    // muestra un mensaje en pantalla
	router.HandleFunc("/index", app.Index)                                // muestra una pagina web utilizando templates
	router.HandleFunc("/sucesos", app.ObtenerSucesos)                     // obtiene el listado de sucesos en JSON
	router.HandleFunc("/suceso/{suceID}", app.ObtenerSuceso)              // obtiene los datos de un suceso en JSON recibiendo por parametro GET el Id de suceso
	router.HandleFunc("/insertar_sala", app.InsertarSala).Methods("POST") // inserta una sala, enviandole los datos a insertar atravez de un metodo POST

	router.HandleFunc("/localidades", app.ObtenerLocalidades)

	//inicio del servicio en el puerto configurado
	http.ListenAndServe(":"+strconv.Itoa(Config.AppConfig.Port), router)
}
