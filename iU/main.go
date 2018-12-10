package main
 import (
	"net/http"
	"github.com/gorilla/mux"
	"log"
	 
 )

 func main(){
	 router := mux.NewRouter()
	 router.HandleFunc("/files", getDocuments).Methods("GET") //Listar
	 router.HandleFunc("/files/{id}", loadDocument).Methods("POST") //Cargar un Doc
	 router.HandleFunc("/files/{id}", removeDocument).Methods("DELETE") // Eliminar un Doc
	//  router.HandleFunc("/users", getUsers).Methods("GET") //Listar usuarios
	//  router.HandleFunc("/users/{id}", loadUsers).Methods("POST") //Cargar un Usuario
	//  router.HandleFunc("/users/{id}", removeUsers).Methods("DELETE") // Eliminar un usuario
	 log.Fatal(http.ListenAndServe(":3000", router))
 }