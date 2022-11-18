package database

import (
	"context"

	"enerbit.com/go/grpc/models"
)

// SetQuestion crea una nueva pregunta en la base de datos
// los parametros son el contexto de la peticion y la pregunta a crear
// si la ejecucion de la consulta es correcta, retorna un error nulo
func (repo *PostgresRepository) SetQuestion(ctx context.Context, question *models.Question) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO questions (id, test_id, question, answer) VALUES ($1, $2, $3, $4)", question.Id, question.TestId, question.Question, question.Answer)
	return err
}
