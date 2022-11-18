package database

import (
	"database/sql"

	// se tiene que agregar con el _ para que el paquete sea agregado a este archivo
	// no se utiliza de forma directa pero es necesario importarlo
	_ "github.com/lib/pq"
)

// Constructor de la estructura PostgresRepository

// para cumplir con la interfaz en ./repository/Repository.go PostgresRepository debe implementar los métodos de la interfaz
// osea, las funciones que se nombran en la interfaz Repository, con sus repectivos parámetros y retornos lo debe
// implementar PostgresRepository
type PostgresRepository struct {
	// db es un puntero a la estructura sql.DB que nos ayuda a conectarnos a la base de datos y ejecutar consultas
	db *sql.DB
}

// Constructor de la estructura PostgresRepository
// recibe como parámetro un string con la conexión a la base de datos
// retorna un puntero a la estructura PostgresRepository para que se pueda utilizar en repository
func NewPostgresRepository(url string) (*PostgresRepository, error) {
	// se conecta a la base de datos y devuelve un puntero a la misma
	db, err := sql.Open("postgres", url)
	// si hay un error, se retorna el error y el puntero a la base de datos es nulo
	if err != nil {
		return nil, err
	}
	// se retorna el puntero a la base de datos y el error es nulo
	return &PostgresRepository{db: db}, nil
}
