package server

import (
	"context"

	"enerbit.com/go/grpc/models"
	"enerbit.com/go/grpc/repository"
	"enerbit.com/go/grpc/studentpb"
)

// constructor del server
// creamos el servidor y le pasamos la interfaz Repository
// el cual tiene la conexión y la implementacion de las consultas a la base de datos
func NewStudentServer(repo repository.Repository) *Server {
	return &Server{repo: repo}
}

// obtenemos lo que nos envia el cliente
// esta es la función que mandamos a crear en el .proto con su respectivo nombre, parametro y retorno
func (s *Server) GetStudent(ctx context.Context, req *studentpb.GetStudentRequest) (*studentpb.Student, error) {
	// req  se puede ver como el message que se creo en el .proto, osea que podemos acceder a los campos
	// hay 2 formas de hacerlo, con el Get y con el campo directamente
	// req.GetId() == req.Id
	// estamos llamando a las funciones que creamos en el repository las cuales son las que se encargan de hacer
	// la implementacion de las funciones que hacen la consulta a la base de datos

	// server => repository => query => database

	student, err := s.repo.GetStudent(ctx, req.GetId())
	// retorna un puntero a nuestro modelo en go
	if err != nil {
		return nil, err
	}

	// se crea un nuevo estudiante con el modelo de nuestro .proto y se le asigna los valores del estudiante que
	// se obtuvo de la base de datos, esto se hace para que el cliente pueda entender el mensaje
	// luego se retorna
	return &studentpb.Student{
		Id:   student.Id,
		Name: student.Name,
		Age:  student.Age,
	}, nil
}

func (s *Server) SetStudent(ctx context.Context, req *studentpb.Student) (*studentpb.SetStudentResponse, error) {
	// aca ya es al reves, nos llega un estudiante con los datos que se quieren guardar
	// convertimos ese dato de nuestro modelo .proto a nuestro modelo en go
	// pasamos ese modelo al repository para que se encargue de hacer la implementacion del query que se realiza a la base de datos

	student := &models.Student{
		Id:   req.GetId(),
		Name: req.GetName(),
		Age:  req.GetAge(),
	}

	err := s.repo.SetStudent(ctx, student)

	if err != nil {
		return nil, err
	}

	return &studentpb.SetStudentResponse{
		Id: student.Id,
	}, nil
}
