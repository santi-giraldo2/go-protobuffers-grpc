package server

import (
	"context"
	"fmt"
	"io"
	"log"

	"enerbit.com/go/grpc/models"
	"enerbit.com/go/grpc/repository"
	"enerbit.com/go/grpc/studentpb"
	"enerbit.com/go/grpc/testpb"
)

type TestServer struct {
	repo repository.Repository
	testpb.UnimplementedTestServiceServer
}

func NewTestServer(repo repository.Repository) *TestServer {
	return &TestServer{repo: repo}
}

func (s *TestServer) GetTest(ctx context.Context, req *testpb.GetTestRequest) (*testpb.Test, error) {
	test, err := s.repo.GetTest(ctx, req.GetId())

	if err != nil {
		return nil, err
	}

	return &testpb.Test{
		Id:   test.Id,
		Name: test.Name,
	}, nil
}

func (s *TestServer) SetTest(ctx context.Context, req *testpb.Test) (*testpb.SetTestResponse, error) {
	test := &models.Test{
		Id:   req.GetId(),
		Name: req.GetName(),
	}

	err := s.repo.SetTest(ctx, test)

	if err != nil {
		return nil, err
	}

	return &testpb.SetTestResponse{
		Id: test.Id,
	}, nil
}

func (s *TestServer) SetQuestions(stream testpb.TestService_SetQuestionsServer) error {
	// tenemos como parametro un stream lo que quiere decir que vamos a obtener n cantidad de datos, osea, varios parametros
	// por eso usamos un for infinito, no tenemos ni idea cuantos van a llegar (que Dios se apiade de nosotros)
	for {
		// Recv es un metodo que nos permite recibir un mensaje del stream
		// asi que lo que hacemos es recibir un mensaje y guardarlo en msg, de a 1
		msg, err := stream.Recv()

		// si el error es EOF, quiere decir que ya no hay mas datos que recibir, quiere decir que el cliente cerró la conexión
		if err == io.EOF {
			// cerramos por nuestra parte el stream y enviamos el tipo de respuesta que se espera en el cliente
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: true,
			})
		}

		// validamos que no haya ocurrido un error al momento de recibir el mensaje, si es asi, mostramos en consola el error
		// y lo retornamos para que el cliente sepa que ocurrió un error
		// con eso no se cierra la conexión del stream, el cliente puede seguir enviando datos
		if err != nil {
			fmt.Println(err)
			return err
		}

		// hacemos ya lo mismo que antes, convertimos el dato, lo mandamos por el repo y validamos que no haya ocurrido un error
		question := &models.Question{
			Id:       msg.GetId(),
			TestId:   msg.GetTestId(),
			Question: msg.GetQuestion(),
			Answer:   msg.GetAnswer(),
		}

		// ahora el contexto esta trabajando en el stream algo asi como en segundo plano, basicamente llega por arte de magia
		err = s.repo.SetQuestion(context.Background(), question)

		// si hay un error al hacer la ejecucion del query mostramos un mensaje en el log del servidor y cerramos la
		// conexión del stream y mandamos el tipo de dato que se espera en la respuesta
		if err != nil {
			fmt.Println(err)
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: false,
			})
		}
	}
}

func (s *TestServer) EnrollStudents(stream testpb.TestService_EnrollStudentsServer) error {
	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: true,
			})
		}

		if err != nil {
			fmt.Println(err)
			return err
		}

		enrollment := &models.Enrollment{
			TestId:    msg.GetTestId(),
			StudentId: msg.GetStudentId(),
		}

		err = s.repo.SetEnrollment(context.Background(), enrollment)

		if err != nil {
			fmt.Println(err)
			return stream.SendAndClose(&testpb.SetQuestionResponse{
				Ok: false,
			})
		}
	}
}

// en este endpoint se recibe un un message y se response con un stream
// asi que para eso necesitamos como parametro la estructura del message y el stream
func (s *TestServer) GetStudentsPerTest(req *testpb.GetStudentsPerTestRequest, stream testpb.TestService_GetStudentsPerTestServer) error {
	// tenemos la conexion con la base de datos que esta en el stream y el id del test por el cual buscar
	// hacemos la consulta el cual nos retorna un slice de estudiantes
	// vamos a mandar cada uno de los estudiantes por el stream de 1 en 1
	students, err := s.repo.GetStudentsPerTest(context.Background(), req.GetTestId())

	if err != nil {
		log.Printf("Getting from repo error: %v", err)
		return err
	}

	// recorremos el slice de estudiantes y vamos mandando cada uno por el stream
	// students == []models.Student
	for key := 0; key < len(students); key++ {
		repoStudent := students[key]
		// convertimos ese model a un message
		student := &studentpb.Student{
			Id:   repoStudent.Id,
			Name: repoStudent.Name,
			Age:  repoStudent.Age,
		}

		// enviamos ese message por el stream
		err := stream.Send(student)

		if err != nil {
			log.Printf("Sending error: %v", err)
			return err
		}
	}

	return nil
}

// en este endpoint se recibe un stream y de devuelve un stream
// solo necesitamos como parametro el stream
func (s *TestServer) TakeTest(stream testpb.TestService_TakeTestServer) error {
	// hacemos el respectivo for para poder estar escuchando cuando el cliente nos envie un mensaje por el stream
	for {
		msg, err := stream.Recv()

		if err == io.EOF {
			return nil
		}

		if err != nil {
			fmt.Println(err)
			return err
		}

		questions, err := s.repo.GetQuestionPerTest(context.Background(), msg.GetTestId())

		if err != nil {
			return err
		}

		var currentQuestion = &models.Question{}
		i := 0
		lenQuestions := len(questions)
		lenQuestions32 := int32(lenQuestions)

		// hacemos un for para poder recorrer el slice que nos retorno el repo/query
		for {
			if i < lenQuestions {
				currentQuestion = questions[i]
				questionToSend := &testpb.QuestionPerTest{
					Id:       currentQuestion.Id,
					Question: currentQuestion.Question,
					Ok:       false,
					Current:  int32(i + 1),
					Total:    lenQuestions32,
				}
				// enviamos la pregunta obtenida por el repo/query por el stream
				err := stream.Send(questionToSend)

				if err != nil {
					return err
				}

				// esperamos a que el cliente nos responda con un mensaje por el stream
				answer, err := stream.Recv()

				if err == io.EOF {
					return nil
				}

				if err != nil {
					return err
				}

				// convertimos la respuesta que nos llego por el stream a un model
				log.Println("Answer: ", answer.GetAnswer())
				answerModel := &models.Answer{
					TestId:     msg.GetTestId(),
					QuestionId: currentQuestion.Id,
					StudentId:  msg.GetStudentId(),
					Answer:     answer.Answer,
					Correct:    (answer.Answer == currentQuestion.Answer),
				}

				// mandamos a guardar esa respuesta en la base de datos
				err = s.repo.SetAnswer(context.Background(), answerModel)

				if err != nil {
					fmt.Println(err)
					return err
				}
				// pasamos de pregunta
				i++
			} else {
				// si no hay ninguna pregunta mandamos el tipo de dato que nos pide el cliente como vacio y cerramos el ciclo
				// ya que si lo dejamos corriendo va a estar enviando mensajes vacios indefinidamente
				// pasamos al siguiente mensaje que nos llega como parametro del cliente
				questionToSend := &testpb.QuestionPerTest{
					Id:       "",
					Question: "",
					Ok:       true,
					Current:  int32(0),
					Total:    int32(0),
				}

				err := stream.Send(questionToSend)

				if err == io.EOF {
					return nil
				}

				if err != nil {
					return err
				}

				break
			}
		}
	}
}

func (s *TestServer) GetTestScore(ctx context.Context, req *testpb.GetTestScoreRequest) (*testpb.TestScore, error) {
	testScore, err := s.repo.GetTestScore(ctx, req.GetTestId(), req.GetStudentId())

	if err != nil {
		return nil, err
	}

	return &testpb.TestScore{
		TestId:    testScore.TestId,
		StudentId: testScore.StudentId,
		Ok:        testScore.Ok,
		Ko:        testScore.Ko,
		Total:     testScore.Total,
		Score:     testScore.Score,
	}, nil
}
