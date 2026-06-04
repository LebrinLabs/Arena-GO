package dados

import "math/rand/v2"


func Rolar(quantidade, lados int) (rolagens []int, total int) {
	for i := 0; i < quantidade; i++{
		valor := rand.IntN(lados) + 1
		rolagens = append(rolagens, valor)
		total += valor
	}
	return rolagens, total
}