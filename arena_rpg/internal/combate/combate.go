package combate

import (
	"arena_rpg/internal/dados"
	"arena_rpg/internal/personagem"
)

type Turno struct {
	Atacante   string `json:"atacante"`
	Defensor   string `json:"defensor"`
	Dano       int    `json:"dano"`
	HPRestante int    `json:"hp_restante"`
}

type Resultado struct {
	Vencedor personagem.Personagem `json:"vencedor"`
	Perdedor personagem.Personagem `json:"perdedor"`
	Turnos   []Turno               `json:"turnos"`
}

func Resolver(a, b personagem.Personagem) Resultado {
	turnos := []Turno{}
	atacante, defensor := a, b

	for atacante.HP > 0 && defensor.HP > 0 {
		_, rolagem := dados.Rolar(1, 6)
		dano := atacante.Ataque + rolagem - defensor.Defesa
		if dano < 1 {
			dano = 1
		}
		defensor.HP -= dano
		turnos = append(turnos, Turno{
			Atacante:   atacante.Nome,
			Defensor:   defensor.Nome,
			Dano:       dano,
			HPRestante: defensor.HP,
		})
		atacante, defensor = defensor, atacante
	}

	if atacante.HP > 0 {
		return Resultado{Vencedor: atacante, Perdedor: defensor, Turnos: turnos}
	}
	return Resultado{Vencedor: defensor, Perdedor: atacante, Turnos: turnos}
}