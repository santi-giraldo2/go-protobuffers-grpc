package models

// nuestro modelo de datos en go
// es decir, la estructura de datos que vamos a utilizar

// el `json:"example"` es para que al momento de serializar/convertir a JSON

type Student struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Age  int32  `json:"age"`
}
