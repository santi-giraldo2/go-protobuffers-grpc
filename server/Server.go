package server

import (
	"enerbit.com/go/grpc/repository"
	"enerbit.com/go/grpc/studentpb"
)

// Server es la estructura que implementa los endpoints de gRPC
type Server struct {
	// con solo eso ya tenemos acceso a las funciones de la interfaz Repository (herencia)
	repo repository.Repository
	// aca es donde se guarda las funciones que creamos en el .proto, en este caso GetStudent y SetStudent
	//
	studentpb.UnimplementedStudentServiceServer
}
