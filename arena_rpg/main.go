package main

import (
	"encoding/json"
	"log"
	"net/http"

	"arena_rpg/internal/personagem"
	"arena_rpg/internal/dados"
)

func main() {

	mux := http.NewServeMux()

	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status":"ok"})
	})

	mux.HandleFunc("POST /personagens", personagem.CriarPersonagem)
	mux.HandleFunc("GET /personagens", personagem.ListarPersonagens)
	mux.HandleFunc("GET /personagens/{id}", personagem.BuscarPersonagemPorId)
	mux.HandleFunc("POST /rolar", dados.RolarHandler)


	log.Println("Servidor rodando na porta 8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

}