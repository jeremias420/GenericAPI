package Config

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb" // Importar el driver para MSSQL
	_ "github.com/go-sql-driver/mysql"   // Importar el driver para MySQL
	"github.com/joho/godotenv"
)

// Variable global para acceder a la configuracion de la app
var AppConfig Configuracion

// Configuracion Estructura que almacena configuracion de la aplicacion
type Configuracion struct {
	DB   DB  `json:"db"`
	Port int `json:"appport"`
	// Host string `json:"host"`
}
type DB struct {
	Driver           string `json:"driver"`
	Server           string `json:"server"`
	User             string `json:"user"`
	Password         string `json:"password"`
	Database         string `json:"database"`
	Port             int    `json:"port"`
	ConnectionString string `json:"connectionString"`
}

func (c Configuracion) IsOk() bool {
	return (c.DB.Driver != "" && c.DB.User != "" && c.DB.Password != "" && c.DB.Database != "")
}

// funcion que retorna la conexion
func ConexionBD() (conexion *sql.DB) {
	var err error
	db, err := sql.Open(AppConfig.DB.Driver, AppConfig.DB.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func LoadConfig() error {
	// Cargar las variables de entorno desde el archivo .properties
	if _, err := os.Stat(".properties"); err == nil {
		godotenv.Load(".properties")
	}

	// host := os.Getenv("APP_HOST")
	// if host == "" {
	// 	host = "0.0.0.0"
	// }

	appPortStr := os.Getenv("PORT")
	if appPortStr == "" {
		appPortStr = os.Getenv("APP_DEFAULT_PORT") // Puerto predeterminado
	}

	// Leer las variables de entorno
	driver := os.Getenv("SQL_DRIVER")
	user := os.Getenv("SQL_USER")
	password := os.Getenv("SQL_PASSWORD")
	server := os.Getenv("SQL_HOST")
	portStr := os.Getenv("SQL_PORT")
	database := os.Getenv("SQL_DATABASE")

	// Convertir el puerto a int
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return errors.New("error al convertir el port")
	}

	// Convertir el puerto a int
	appPort, err := strconv.Atoi(appPortStr)
	if err != nil {
		return errors.New("error al convertir el appPort")
	}

	// Crear la cadena de conexión dependiendo del driver
	var connectionString string
	switch driver {
	case "mysql":
		connectionString = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, server, port, database)
	case "mssql":
		connectionString = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%d;database=%s;encrypt=%s;", server, user, password, port, database, "disable")
		//connectionString = fmt.Sprintf("sqlserver://%s:%s@%s:%d",config.DB.User, config.DB.Password, config.DB.Server, config.DB.Port)
	default:
		return errors.New("driver no soportado")
	}

	// Crear la estructura de configuración
	AppConfig = Configuracion{
		DB: DB{
			Driver:           driver,
			Server:           server,
			User:             user,
			Password:         password,
			Port:             port,
			Database:         database,
			ConnectionString: connectionString,
		},
		Port: appPort,
		// Host: host,
	}

	return nil

}
