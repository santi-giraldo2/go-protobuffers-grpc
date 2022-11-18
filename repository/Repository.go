package repository

import (
	"context"

	"enerbit.com/go/grpc/models"
)

// una interfaz es un conjunto de funciones que se deben implementar en una estructura
// en este caso, la interfaz Repository es un conjunto de funciones que se deben implementar en la estructura PostgresRepository
// osea, todas las funciones que estan nombradas, junto a sus parametros y retorno, deben estar en la estructura PostgresRepository
// la estructura PostgresRepository es la que se crea en database\postgres.go para manejar la base de datos
type Repository interface {
	GetStudent(ctx context.Context, id string) (*models.Student, error)
	SetStudent(ctx context.Context, student *models.Student) error
	GetTest(ctx context.Context, id string) (*models.Test, error)
	SetTest(ctx context.Context, test *models.Test) error
	SetQuestion(ctx context.Context, question *models.Question) error
	SetEnrollment(ctx context.Context, enrollment *models.Enrollment) error
	GetStudentsPerTest(ctx context.Context, testId string) ([]*models.Student, error)
	GetQuestionPerTest(ctx context.Context, testId string) ([]*models.Question, error)
	SetAnswer(ctx context.Context, answer *models.Answer) error
	GetTestScore(ctx context.Context, testId, studentId string) (*models.TestScore, error)
}

// variable que se usa para almacenar la implementacion de la interfaz Repository
var implementation Repository

func SetRepository(repository Repository) {
	implementation = repository
}

// se hace la implementacion de las funciones de la interfaz Repository
// cada función que se nombra en la interfaz Repository se crean en esta parte, con los mismos parametros y retorno
// lo que hacemos es hacer un intermediario, cuando se quiera hacer una consulta a la base de datos hacemos el llamado
//  a la función de la interfaz Repository el cual se encarga de llamar a la función de la estructura de la DB (PostgresRepository)
