Este es un proyecto para el entendimiento de como crear un servidor y un cliente en Golang junto con la tecnología de gRPC

para correr el proyecto es necesario crear el contenedor que se encuentra en database

docker build . -t grpc-db
docker run -p 54321:5432 grpc-db

para correr el servidor de student
go run .\server-student\main.go

para correr el servidor de test
go run .\server-test\main.go

para correr el cliente en gRPC
go run .\client\main.go

para comenzar con un proyecto de gRPC en Golang es necesario iniciar creando los archivos .proto
en un archivo .proto se define la estructura de los datos que se van a enviar y recibir en el servidor y el cliente.
para definir que tipo de dato se va a enviar se utiliza message el cual se define de la siguiente manera
message Student {
int32 id = 1;
string name = 2;
string email = 3;
string phone = 4;
string address = 5;
}

basicamente es como un JSON, se define que tipo de dato es y el nombre de la variable, el numero que se encuentra al lado del nombre de la variable es el numero de la variable que se va a enviar, es decir, si se envia un JSON con 5 variables, la primera variable se le asigna el numero 1, la segunda el numero 2 y asi sucesivamente.
el tipo de dato tambien puede ser otro message

para definir el endpoint se utiliza service, se le tiene que indicar el nombre que va a tener, que parametros va a recibir y que tipo de dato va a retornar:

service StudentService {
rpc GetStudent (StudentRequest) returns (StudentResponse) {}
}

en este caso se define que el endpoint se va a llamar GetStudent, que va a recibir un StudentRequest y va a retornar un StudentResponse

una vez realizado el archivo .proto se tiene que generar el codigo de Golang para poder utilizarlo en el servidor y el cliente, para esto se utiliza el siguiente comando:

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative <nombre del archivo proto>.proto

este comando genera dos archivos, uno con la extension .pb.go y otro con la extension .pb.gw.go, el primero es el archivo que se va a utilizar en el servidor y el cliente, el segundo es el archivo que se utiliza para la comunicacion con el servidor.
el comando anterior tiene unos flags para que el codigo generado se cree en la misma carpeta que el archivo .proto

para iniciar un proyecto Golang se utiliza el siguiente comando:
go mod init <nombre del proyecto>

este comando crea un archivo go.mod en el cual se especifica el nombre del proyecto y las dependencias que se van a utilizar en el proyecto

para instalar las dependencias se utiliza el siguiente comando:
go get <nombre de la dependencia>

los paquetes que se utilizan en este proyecto son:
google.golang.org/grpc => para la comunicacion gRPC
google.golang.org/protobuf => para generar el codigo de los archivos .proto
github.com/lib/pq => paquete para la conexion con la base de datos postgres

se puede utilizar el comando
go mod tidy
para que se actualice el archivo go.mod con las dependencias que se estan utilizando en el proyecto y se eliminen las que no se estan utilizando.

explicacion de las carpetas:

client:
en esta carpeta se encuentra el cliente de gRPC, en el archivo main.go se encuentra la conexion con el servidor y las llamadas a los endpoints

database:
contiene el archivo Dockerfile para crear el contenedor de la base de datos
PostgresRepository.go tiene la conexion con la base de datos y las funciones para realizar las consultas a la base de datos

models:
contiene los modelos en Golang de los datos que se van a enviar y recibir en el servidor y el cliente
se le agrega el tag json para que al momento de enviar los datos en formato JSON se pueda enviar con el nombre de la key JSON y no con el nombre de la variable en Golang

repository:
realiza la impletación de los metodos que se encuentran en el archivo database/PostgresRepository.go
es el intermediario entre los endpoints y los querys de la base de datos

server:
aqui se indica como se va a comportar el servidor, cuales endpoints va a tener y que funciones se van a ejecutar cuando se llame a un endpoint

server-student:
se crea un servidor, se le indica que puerto va a escuchar y se le indica que endpoints va a tener

server-test:
se crea un servidor, se le indica que puerto va a escuchar y se le indica que endpoints va a tener

studentpb:
contiene los archivos .pb.go y .pb.gw.go que se generan con el comando protoc y el archivo .proto

testpb:
contiene los archivos .pb.go y .pb.gw.go que se generan con el comando protoc y el archivo .proto
