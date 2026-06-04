package personagem

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type Personagem struct {
	ID     int    `json:"id"`
	Nome   string `json:"nome"`
	Classe string `json:"classe"`
	HP     int    `json:"hp"`
	Ataque int    `json:"ataque"`
	Defesa int    `json:"defesa"`
	Nivel  int    `json:"nivel"`
	XP     int    `json:"xp"`
}

var db *sql.DB

func IniciarDB(caminho string) error {
	var err error
	db, err = sql.Open("sqlite", caminho)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS personagens (
			id     INTEGER PRIMARY KEY AUTOINCREMENT,
			nome   TEXT NOT NULL,
			classe TEXT,
			hp     INTEGER,
			ataque INTEGER,
			defesa INTEGER,
			nivel  INTEGER,
			xp     INTEGER
		)
	`)
	return err
}

func Inserir(p Personagem) (Personagem, error) {
	res, err := db.Exec(
		`INSERT INTO personagens (nome, classe, hp, ataque, defesa, nivel, xp)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		p.Nome, p.Classe, p.HP, p.Ataque, p.Defesa, p.Nivel, p.XP,
	)
	if err != nil {
		return Personagem{}, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return Personagem{}, err
	}
	p.ID = int(id)
	return p, nil
}

func Listar() ([]Personagem, error) {
	rows, err := db.Query(`SELECT id, nome, classe, hp, ataque, defesa, nivel, xp FROM personagens`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lista := []Personagem{}
	for rows.Next() {
		var p Personagem
		if err := rows.Scan(&p.ID, &p.Nome, &p.Classe, &p.HP, &p.Ataque, &p.Defesa, &p.Nivel, &p.XP); err != nil {
			return nil, err
		}
		lista = append(lista, p)
	}
	return lista, rows.Err()
}

func Buscar(id int) (Personagem, bool) {
	var p Personagem
	err := db.QueryRow(
		`SELECT id, nome, classe, hp, ataque, defesa, nivel, xp FROM personagens WHERE id = ?`,
		id,
	).Scan(&p.ID, &p.Nome, &p.Classe, &p.HP, &p.Ataque, &p.Defesa, &p.Nivel, &p.XP)
	if err != nil {
		return Personagem{}, false // não achou (ou erro)
	}
	return p, true
}

func (p *Personagem) GanharXP(xp int) {
	p.XP += xp
	for p.XP >= 100 {
		p.XP -= 100
		p.Nivel++
		p.HP += 10
		p.Ataque += 2
		p.Defesa += 1
	}
}

func DarXP(id, xp int) (Personagem, bool) {
	p, ok := Buscar(id)
	if !ok {
		return Personagem{}, false
	}
	p.GanharXP(xp)
	_, err := db.Exec(
		`UPDATE personagens SET hp = ?, ataque = ?, defesa = ?, nivel = ?, xp = ? WHERE id = ?`,
		p.HP, p.Ataque, p.Defesa, p.Nivel, p.XP, p.ID,
	)
	if err != nil {
		return Personagem{}, false
	}
	return p, true
}