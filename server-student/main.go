package main

import (
	"log"
	"net"

	"enerbit.com/go/grpc/database"
	"enerbit.com/go/grpc/server"
	"enerbit.com/go/grpc/studentpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// se indica en que puerto va a correr el servidor
	list, err := net.Listen("tcp", ":5060")

	// si hay un error se imprime y se sale del programa
	if err != nil {
		log.Fatal(err)
	}

	// defer: se ejecuta cuando la funcion termina
	// cerramos la conexion al puerto
	defer list.Close()
	// se abre la conexion a la base de datos
	repo, err := database.NewPostgresRepository("postgres://postgres:postgres@localhost:54321/postgres?sslmode=disable")

	// si hay un error se imprime y se sale del programa
	if err != nil {
		log.Fatal(err)
	}

	// creamos el servidor
	server := server.NewStudentServer(repo)
	// se crea el servidor gRPC
	s := grpc.NewServer()

	// combinamos el servidor que hemos destinado para ser un gRPC
	// el servidor que tiene la implementación a las funciones que interactuan con la base de datos
	// y el servidor que se encarga de la comunicación gRPC
	studentpb.RegisterStudentServiceServer(s, server)
	// registramos ese sevidor gRPC para que puede ser usado/lo publicamos
	reflection.Register(s)

	// se inicia el servidor gRPC y se le pasa el puerto
	if err := s.Serve(list); err != nil {
		log.Fatal(err)
	}
}
