package database

import (
	"context"
	"log"

	"enerbit.com/go/grpc/models"
)

// SetTest crea un nuevo test en la base de datos
// los parametros son el contexto de la peticion y el test a crear
// si la ejecucion de la consulta es correcta, retorna un error nulo
func (repo *PostgresRepository) SetTest(ctx context.Context, test *models.Test) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO tests (id, name) VALUES ($1, $2)", test.Id, test.Name)
	return err
}

// GetTest obtiene un test de la base de datos
// los parametros son el contexto de la peticion y el id del test
// si la ejecucion de la consulta es correcta, retorna un error nulo
func (repo *PostgresRepository) GetTest(ctx context.Context, id string) (*models.Test, error) {
	// hace la misma consulta que GetStudent
	// obtiene el test de la base de datos
	// se cierra la conexion a la base de datos
	// se crea un test de nuestro modelo de datos
	// se recorre cada fila de la consulta
	// se asigna el valor de cada columna a las propiedades del test
	// se retorna el test y el error es nulo
	rows, err := repo.db.QueryContext(ctx, "SELECT id, name FROM tests WHERE id = $1", id)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var test = &models.Test{}
	for rows.Next() {
		err := rows.Scan(&test.Id, &test.Name)

		if err != nil {
			return nil, err
		}
	}

	return test, nil
}

// Obtiene los estudiantes que estan inscritos en un test
func (repo *PostgresRepository) GetStudentsPerTest(ctx context.Context, testId string) ([]*models.Student, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, name, age FROM students WHERE id IN (SELECT student_id FROM enrollments WHERE test_id = $1)", testId)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()

		if err != nil {
			log.Fatal(err)
		}
	}()

	var students []*models.Student
	for rows.Next() {
		var student = models.Student{}

		if err := rows.Scan(&student.Id, &student.Name, &student.Age); err == nil {
			students = append(students, &student)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return students, nil
}

// inscribe a un estudiante en un test
func (repo *PostgresRepository) SetEnrollment(ctx context.Context, enrollment *models.Enrollment) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO enrollments (test_id, student_id) VALUES ($1, $2)", enrollment.TestId, enrollment.StudentId)
	return err
}

// obtiene los resultados de un test
func (repo *PostgresRepository) GetQuestionPerTest(ctx context.Context, testId string) ([]*models.Question, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, question, answer FROM questions WHERE test_id = $1", testId)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()

		if err != nil {
			log.Fatal(err)
		}
	}()

	var questions []*models.Question
	for rows.Next() {
		var question = models.Question{}

		if err := rows.Scan(&question.Id, &question.Question, &question.Answer); err == nil {
			questions = append(questions, &question)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return questions, nil
}

// agrega una respuesta a un test
func (repo *PostgresRepository) SetAnswer(ctx context.Context, answer *models.Answer) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO answers (student_id, test_id, question_id, answer, correct) VALUES($1, $2, $3, $4, $5)", answer.StudentId, answer.TestId, answer.QuestionId, answer.Answer, answer.Correct)
	return err
}

// Obtiene el promedio de un estudiante respecto a un test
func (repo *PostgresRepository) GetTestScore(ctx context.Context, testId, studentId string) (*models.TestScore, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT correct FROM answers WHERE test_id = $1 AND student_id = $2", testId, studentId)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()

		if err != nil {
			log.Fatal(err)
		}
	}()

	var answer = models.Answer{}
	var testScore = models.TestScore{
		TestId:    testId,
		StudentId: studentId,
	}

	for rows.Next() {
		err := rows.Scan(&answer.Correct)

		if err != nil {
			return nil, err
		}

		testScore.Total += 1

		if answer.Correct {
			testScore.Ok += 1
		} else {
			testScore.Ko += 1
		}
	}

	testScore.Score = testScore.Ok * 10 / testScore.Total

	return &testScore, nil
}
