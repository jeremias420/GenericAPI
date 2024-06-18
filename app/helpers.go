package app

import (
	Config "GenericAPI/app/config"
	"encoding/json"
	"io"
	"net/http"
)

// var plantillas = template.Must(template.ParseGlob("templates/*"))

type Option struct {
	ID          string
	Description string
}

type PageData struct {
	ComboOptions []Option
}

// Define una estructura para el resultado JSON.
type ResultadoJSON struct {
	Data []map[string]interface{} `json:"data"`
}

// Función que ejecuta una consulta SQL y devuelve el resultado en formato JSON.
func ConsultarSQL(query string) (string, error) {
	// Obtén una conexión a la base de datos.
	db := Config.ConexionBD()
	defer db.Close()

	// Ejecuta la consulta SQL.
	rows, err := db.Query(query)
	if err != nil {
		// log.Println("Error al ejecutar la consulta SQL:", err)
		return "", err
	}
	defer rows.Close()

	// Itera sobre los resultados y construye un arreglo de mapas.
	var result []map[string]interface{}
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return "", err
		}

		entry := make(map[string]interface{})
		for i, col := range columns {
			switch v := values[i].(type) {
			case nil:
				entry[col] = nil
			case []byte:
				entry[col] = string(v)
			default:
				entry[col] = v
			}
		}
		result = append(result, entry)
	}

	// Convierte el resultado a JSON.
	jsonData, err := json.Marshal(ResultadoJSON{Data: result})
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func ObtenerConsultaJSON(w http.ResponseWriter, r *http.Request, consultaSQL string) {
	// Ejecutar la consulta SQL
	resultado, err := ConsultarSQL(consultaSQL)
	if err != nil {
		http.Error(w, "Error al obtener sucesos: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Configurar la cabecera HTTP para indicar que la respuesta es JSON
	w.Header().Set("Content-Type", "application/json")

	// Escribir el resultado en la respuesta HTTP
	_, err = w.Write([]byte(resultado))
	if err != nil {
		http.Error(w, "Error al escribir la respuesta: "+err.Error(), http.StatusInternalServerError)
		return
	}

}
func ObtenerParametroPOST(body []byte, name string) string {

	// Declara un mapa para almacenar los datos del JSON.
	var jsonData map[string]interface{}

	// Deserializa el JSON en el mapa.
	err := json.Unmarshal(body, &jsonData)
	if err != nil {
		return ""
	}

	value, existe := jsonData[name]
	if existe {
		if strValue, ok := value.(string); ok {
			return strValue
		}
	}

	return ""
}

func Comillas(value string) string {
	return "'" + value + "'"
}

func GetBody(r *http.Request) []byte {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil
	}
	return body
}
