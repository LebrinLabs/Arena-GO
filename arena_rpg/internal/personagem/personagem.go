package personagem

import "sync"

type Personagem struct {
	Id 		int 	`json:"id"`
	Nome 	string	`json:"nome"`
	Raça 	string	`json:"raça"`
	Classe 	string	`json:"classe"`
	HP 		int		`json:"hp"`
	Ataque 	int		`json:"ataque"`
	Defesa 	int		`json:"defesa"`
	Nivel 	int		`json:"nivel"`
	Xp 		int		`json:"xp"`
	Armas 	[]string `json:"armas"`
}


var ( 

	personagens = make(map[int]Personagem)
	proximoId = 1
	mu sync.Mutex

)

func Buscar(id int) (Personagem, bool) {
	mu.Lock()
	p, ok := personagens[id]
	mu.Unlock()
	return p, ok
}

func (p *Personagem) GanharXP(xp int) {
	p.Xp += xp
	for p.Xp >= 100{
		p.Xp -= 100
		p.Nivel++
		p.HP += 10
		p.Ataque += 2
		p.Defesa += 1
	}
}

func DarXP(id int, xp int) (Personagem, bool) {
	mu.Lock()
	p, ok := personagens[id]
	if ok {
		p.GanharXP(xp)
		personagens[id] = p
	}
	mu.Unlock()
	return p, ok
}