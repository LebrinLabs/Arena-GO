package personagem

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func CriarPersonagem(w http.ResponseWriter, r *http.Request) {
	var p Personagem

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, "dados inválidos", http.StatusBadRequest)
		return
	}

	mu.Lock()
	p.Id = proximoId
	proximoId++
	p.Nivel = 1
	p.Xp = 0
	personagens[p.Id] = p
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 Created
	json.NewEncoder(w).Encode(p)

}

func ListarPersonagens(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	lista := make([]Personagem, 0, len(personagens))

	for _, p := range personagens {
		lista = append(lista, p)
	}
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lista)

}

func BuscarPersonagemPorId( w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "id inválido", http.StatusBadRequest)
		return
	}

	mu.Lock()
	p, ok := personagens[id]
	mu.Unlock()
	
	if !ok {
		http.Error(w, "personagem não encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}