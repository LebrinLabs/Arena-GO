# Criando Personaem

type Personagem struct {
	Id 		int 	`json:"id"`
	Nome 	string	`json:"nome"`
	Raça 	string	`json:"raça"`
	Classe 	string	`json:"classe"`
	Hp 		int		`json:"hp"`
	Ataque 	int		`json:"ataque"`
	Defesa 	int		`json:"defesa"`
	Nivel 	int		`json:"nivel"`
	Xp 		int		`json:"xp"`
	Armas 	[]string `json:"armas"`
}

var personagens = make(map[int]Personagem)
var proximoId = 1

