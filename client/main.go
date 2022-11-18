package main

import (
	"context"
	"io"
	"log"

	"enerbit.com/go/grpc/testpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// se crea un cliente de gRPC
	cc, err := grpc.Dial("localhost:5070", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	// se cierra la conexión cuando se termine de ejecutar el programa
	defer cc.Close()

	// se obtiene una instancia del cliente de gRPC
	// esa funcion es creada de forma automatica por el compilador de protobuf
	c := testpb.NewTestServiceClient(cc)

	// estas funciones es un ejemplo de las 4 formas de comunicación que se pueden hacer con gRPC
	DoUnary(c)
	//DoClientStreaming(c)
	//DoServerStreaming(c)
	// DoBidireccionalStreaming(c)
}

func DoUnary(c testpb.TestServiceClient) {
	req := &testpb.GetTestRequest{
		Id: "t1",
	}

	// aunque la funcion GetTest pareciera que solo se puede utilizar en el service de hecho tambien lo podemos
	// utilizar en el cliente, esto es posible gracias a que el compilador de protobuf genera la funcion en ambos lados

	// asi que si hacemos un Get en el .proto, podemos utilizar la mismo llamado al endpoint en el cliente y en el server
	res, err := c.GetTest(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling GetTest: %v", err)
	}

	log.Printf("GetTest response: %v", res)
}

// el cliente envia un stream de datos al servidor y el servidor responde con un solo mensaje
func DoClientStreaming(c testpb.TestServiceClient) {
	// se crea los mensajes que se van a enviar al servidor
	questions := []*testpb.Question{
		{
			Id:       "q8t1",
			Answer:   "lijasd",
			Question: "lijasdsad",
			TestId:   "t1",
		},
		{
			Id:       "q9t1",
			Answer:   "owiqeu",
			Question: "oiwque",
			TestId:   "t1",
		},
		{
			Id:       "q10t1",
			Answer:   "mnsd",
			Question: "okkwqwwwwww",
			TestId:   "t1",
		},
	}

	// se conecta al endpoint del servidor
	stream, err := c.SetQuestions(context.Background())

	if err != nil {
		log.Fatalf("Error while calling SetQuestions: %v", err)
	}

	// se recorre el slice de mensajes
	for _, question := range questions {
		log.Println("Sending question: ", question.GetId())
		// se envia el mensaje al servidor
		stream.Send(question)
	}

	// se cierra el stream y se obtiene la respuesta del servidor
	msg, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error while closing stream: %v", err)
	}

	// se imprime en consola la respuesta del servidor
	log.Printf("response from server: %v", msg)
}

// el cliente envia un mensaje al servidor y el servidor responde con un stream de datos
func DoServerStreaming(c testpb.TestServiceClient) {
	// se crea el mensaje que se va a enviar al servidor
	req := &testpb.GetStudentsPerTestRequest{
		TestId: "t1",
	}

	// se manda el mensaje al servidor y se obtiene un stream de datos
	stream, err := c.GetStudentsPerTest(context.Background(), req)

	if err != nil {
		log.Fatalf("Error while calling GetStudentsPerTest: %v", err)
	}

	// se hace un bucle infinito para recibir los mensajes del servidor
	for {
		// se recibe el mensaje del servidor
		msg, err := stream.Recv()

		// Se verifica si la conexión al stream del servidor se cerro y se sale del bucle
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}

		// se imprime en consola el mensaje del servidor
		log.Printf("Server response: %v", msg)
	}

	// se cierra el stream de parte del cliente
	closeErr := stream.CloseSend()

	if closeErr != nil {
		log.Fatalf("Error while closing stream: %v", closeErr)
	}

	log.Println("Connection closed")
}

// el cliente envia un stream de datos al servidor y el servidor responde con un stream de datos
func DoBidireccionalStreaming(c testpb.TestServiceClient) {
	startTest := &testpb.TakeTestRequest{
		TestId: "t1",
	}
	testAnswer := &testpb.TakeTestRequest{
		TestId: "t1",
		Answer: "asdasdas",
	}
	// se conecta al endpoint del servidor y se obtiene un stream de datos
	stream, err := c.TakeTest(context.Background())

	if err != nil {
		log.Fatalf("Error while calling TakeTest: %v", err)
	}

	// se envia el primer mensaje al servidor
	stream.Send(startTest)

	for {
		// se recibe el mensaje del servidor
		msg, err := stream.Recv()

		// Cuando el servidor cierra la conexión se sale del bucle
		// o cuando el server response con un mensaje true
		if err == io.EOF || msg.GetOk() {
			break
		}

		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}

		log.Printf("question %v", msg)
		// se envia el mensaje al servidor
		stream.Send(testAnswer)
	}

	// se cierra el stream de parte del cliente
	closeErr := stream.CloseSend()

	if closeErr != nil {
		log.Fatalf("Error while closing stream: %v", closeErr)
	}

	log.Println("Connection closed")
}
