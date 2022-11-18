package database

import (
	"context"
	"log"

	"enerbit.com/go/grpc/models"
)

// SetStudent crea un nuevo estudiante en la base de datos
// los parametros son el contexto de la petición y el estudiante a crear
// si la ejecucion de la consulta es correcta, retorna un error nulo
func (repo *PostgresRepository) SetStudent(ctx context.Context, student *models.Student) error {
	// no queremos saber que retorna la funcion, solo queremos saber si hubo un error para continuar
	_, err := repo.db.ExecContext(ctx, "INSERT INTO students (id, name, age) VALUES ($1, $2, $3)", student.Id, student.Name, student.Age)
	return err
}

// GetStudent obtiene un estudiante de la base de datos
// los parametros son el contexto de la peticion y el id del estudiante
// si la ejecucion de la consulta es correcta, retorna un error nulo
func (p *PostgresRepository) GetStudent(ctx context.Context, id string) (*models.Student, error) {
	// ejecuta la consulta en la base de datos y retorna un PUNTERO a las FILAS
	// estas filas tienen el modelo de la base de datos. No es el mismo al que tenemos en la carpeta models
	rows, err := p.db.QueryContext(ctx, "SELECT id, name, age FROM students WHERE id = $1", id)
	// validacion y retorno de error
	if err != nil {
		return nil, err
	}
	// defer => se ejecuta al final de la funcion
	defer func() {
		// se cierra la conexion a la base de datos
		err := rows.Close()
		// validacion del error y detencion del programa
		if err != nil {
			log.Fatal(err)
		}
	}() // es una funcion anonima

	// se crea un estudiante de nuestro modelo de datos
	// se guarda en memoria
	student := &models.Student{}

	// se recorre cada fila de la consulta
	for rows.Next() {
		// se asigna el valor de cada columna a las propiedades del estudiante
		// es decir
		// go.id = db.Id
		// go.name = db.Name
		// go.age = db.Age
		err := rows.Scan(&student.Id, &student.Name, &student.Age)
		// validacion y retorn de la funcion GetStudent
		if err != nil {
			return nil, err
		}
	}
	// se retorna el estudiante y el error es nulo
	return student, nil
}
