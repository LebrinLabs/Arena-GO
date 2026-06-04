package combate

import (
	"encoding/json"
	"net/http"

	"arena_rpg/internal/personagem"
)

func CombateHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AtacanteId int `json:"atacante_id"`
		DefensorId int `json:"defensor_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "dados inválidos", http.StatusBadRequest)
		return
	}

	atacante, ok1 := personagem.Buscar(req.AtacanteId)
	defensor, ok2 := personagem.Buscar(req.DefensorId)

	if !ok1 || !ok2 {
		http.Error(w, "personagem não encontrado", http.StatusNotFound)
		return
	}

	resultado := Resolver(atacante, defensor)

	personagem.DarXP(resultado.Vencedor.ID, 50)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resultado)
}