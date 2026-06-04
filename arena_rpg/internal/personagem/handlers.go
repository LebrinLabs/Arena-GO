package main

import (
	"encoding/json"
	"net/http"
)

func criarPersonagem(w http.ResponseWriter, r *http.Request) {
	var p Personagem

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "dados inválidos", http.StatusBadRequest)
		return
	}

	mu.Lock()
	p.ID = proximoId
	proximoId++
	p.Nivel = 1
	p.Xp = 0
	personagens[p.ID] = p
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created
	json.NewEncoder(w).Encode(p)

}

func listarPersonagens(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	lista := make([]Personagem, 0, Len(personagens))

	for _, p := range personagens {
		lista = append(lista, p)
	}
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lista)

}