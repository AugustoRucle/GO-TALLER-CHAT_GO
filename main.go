package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Response struct {
	Message string `json: "message"`
	Status  int    `json: "status"`
	IsValid bool   `json: "isvalid"`
}

//Se almacena toda la estructura
var Users = struct {
	m            map[string]User //Se almacen todos los usuario que se conecten
	sync.RWMutex                 //Para evitar que el diccionario se rompa en el momento de haber goruntimes
}{m: make(map[string]User)}

type User struct { // ---> ¿Que sucede?
	User_Name string
	WebSocket *websocket.Conn
}

func HolaMundo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hola mundo desde Go"))
}

func HolaMundoJson(w http.ResponseWriter, r *http.Request) {
	response := CreateResponse("Hola mundo increible", 1, true)
	json.NewEncoder(w).Encode(response)
}

func CreateResponse(message string, status int, valid bool) Response {
	return Response{message, status, valid}
}

func LoadStatic(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "HTML/index.html")
}

func UserExit(user_name string) bool {
	Users.RLock() //---> ¿Cómo funciona?
	defer Users.RUnlock()
	if _, ok := Users.m[user_name]; ok {
		return true
	}
	return false
}

func Validate(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // --> ¿Qué es parse?
	user_name := r.FormValue("user_name")
	response := Response{}

	if UserExit(user_name) {
		response.IsValid = false
	} else {
		response.IsValid = true
	}
	json.NewEncoder(w).Encode(response)
}

func CreateUser(user_name string, ws *websocket.Conn) User {
	return User{user_name, ws}
}

func AddUser(user User) {
	Users.Lock()
	defer Users.Unlock()

	Users.m[user.User_Name] = user

}

func RemoveUser(user_name string) {
	Users.Lock()
	defer Users.Unlock()
	delete(Users.m, user_name)
}

func SendMessage(type_mesage int, message []byte) {
	Users.RLock()
	defer Users.RUnlock()

	for _, user := range Users.m {
		if err := user.WebSocket.WriteMessage(type_mesage, message); err != nil {
			return
		}
	}
}

func ToArrayByte(value string) []byte {
	return []byte(value)
}

func ConactMessage(user_name string, arreglo []byte) string {
	return user_name + ":" + string(arreglo[:])
}

func WebSocket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)            //Obtenemos los datos de la url
	user_name := vars["user_name"] //Obtenemos el user_name dentro de la url
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if err != nil {
		log.Println(err)
		return
	}
	current_user := CreateUser(user_name, ws)
	AddUser(current_user)
	log.Println("Nuevo usuario agregado")

	for {
		type_message, message, err := ws.ReadMessage()
		if err != nil {
			RemoveUser(user_name)
			return
		}
		final_message := ConactMessage(user_name, message)
		log.Println("Data: ", final_message)
		SendMessage(type_message, ToArrayByte(final_message))
	}
}

func main() {
	cssHandle := http.FileServer(http.Dir("CSS/"))
	jsHandle := http.FileServer(http.Dir("JS/"))
	jqueryHandle := http.FileServer(http.Dir("JQUERY/"))

	mux := mux.NewRouter()
	mux.HandleFunc("/Hola", HolaMundo).Methods("GET")
	mux.HandleFunc("/HolaJson", HolaMundoJson).Methods("GET")
	mux.HandleFunc("/index", LoadStatic).Methods("GET")
	mux.HandleFunc("/validate", Validate).Methods("POST")
	mux.HandleFunc("/chat/{user_name}", WebSocket).Methods("GET")

	http.Handle("/", mux)
	http.Handle("/CSS/", http.StripPrefix("/CSS/", cssHandle))
	http.Handle("/JS/", http.StripPrefix("/JS/", jsHandle))
	http.Handle("/JQUERY/", http.StripPrefix("/JQUERY/", jqueryHandle))

	log.Println("Servidor Activo, puerto[8000]")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
