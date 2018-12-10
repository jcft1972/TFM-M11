package main
 import (
	"net/http"
	"github.com/gorilla/mux"
	"log"
	"fmt"
	"encoding/json"

	 
 )

 func main(){
	 router := mux.NewRouter()
	 router.HandleFunc("/files", getDocuments).Methods("GET") //Listar
	 router.HandleFunc("/files/{id}", loadDocument).Methods("POST") //Cargar un Doc
	 router.HandleFunc("/files/{id}", removeDocument).Methods("DELETE") // Eliminar un Doc
	 router.HandleFunc("/users", getUsers).Methods("GET") //Listar usuarios
	 router.HandleFunc("/users/{name} {email}", createUser).Methods("POST") //Cargar un Usuario
	//  router.HandleFunc("/users/{id}", removeUsers).Methods("DELETE") // Eliminar un usuario
	 log.Fatal(http.ListenAndServe(":3000", router))
 }

 func getDocuments(w http.ResponseWriter, r *http.Request) {
	var ListFiles string
	flag	:= true
	prepareListFiles()
	for flag {
	 flag , ListFiles = getListFiles()
	 fmt.Println(flag, ListFiles)
	}
	//json.NewEncoder(w).Encode(ListFiles)
	fmt.Fprintf(w, ListFiles)
}

func loadDocument(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)	
	fmt.Println(params["id"])
	infoFile := loadFile(params["id"])
	//infoFile := loadFile("C:\\Users\\jfernandez\\Pictures\\LogoFLSA_135.jpg")
	sendMails()
	json.NewEncoder(w).Encode(infoFile)
	
}

func removeDocument(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)	
	removeFile(params["id"])
	json.NewEncoder(w).Encode("Solicitud enviada.")

}

func getUsers(w http.ResponseWriter, r *http.Request) {
	users := readList()
	//json.NewEncoder(w).Encode(users)
	fmt.Fprint(w, users)
}

func createUser(w http.ResponseWriter, r *http.Request){
	params := mux.Vars(r)	
	writeList(createMail(readList(), params["name"], params["email"]))
}
