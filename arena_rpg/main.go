package main

import (
	"encoding/json"
	"log"
	"net/http"

	"arena_rpg/internal/personagem"

)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status":"ok"})
	})

	mux.HandleFunc("POST /personagens", personagem.CriarPersonagem)
	mux.HandleFunc("GET /personagens", personagem.ListarPersonagens)

	log.Println("Servidor rodando na porta 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

}