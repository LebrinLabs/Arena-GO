package dados

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type respostaRolagem struct {
	Dados 		string 	`json:"dados"`
	Rolagens 	[]int 	`json:"rolagens"`
	Total 		int 	`json:"total"`
}

func RolarHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Dados string `json:"dados"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "dados inválidos", http.StatusBadRequest)
		return
	}

	qtdStr, ladosStr, ok := strings.Cut(req.Dados, "d")
	if !ok {
		http.Error(w, "formato inválido, use algo com 3d6", http.StatusBadRequest)
		return
	}

	qtd, err1 := strconv.Atoi(qtdStr)
	lados, err2 := strconv.Atoi(ladosStr)

	if err1 != nil || err2 != nil || qtd <= 0 || lados <= 1 {
		http.Error(w, "quantidade e lados devem ser números positivos", http.StatusBadRequest)
		return
	}

	if qtd > 100 {
		http.Error(w, "quantidade muito grande, limite de 100", http.StatusBadRequest)
		return
	}

	rolagens, total := Rolar(qtd, lados)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respostaRolagem{
		Dados:    req.Dados,
		Rolagens: rolagens,
		Total:    total,
	})
}
