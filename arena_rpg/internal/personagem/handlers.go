package personagem

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func CriarPersonagem(w http.ResponseWriter, r *http.Request) {
	var p Personagem
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}
	p.Nivel = 1
	p.XP = 0

	p, err := Inserir(p)
	if err != nil {
		http.Error(w, "erro ao salvar", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func ListarPersonagens(w http.ResponseWriter, r *http.Request) {
	lista, err := Listar()
	if err != nil {
		http.Error(w, "erro ao listar", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lista)
}

func BuscarPersonagemPorId(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "id inválido", http.StatusBadRequest)
		return
	}

	p, ok := Buscar(id)
	if !ok {
		http.Error(w, "personagem não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}