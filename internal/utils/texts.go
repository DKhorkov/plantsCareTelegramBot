package utils

import (
	"fmt"
)

// GetWateringInterval - отдает тектосове выражения интервала полива.
// “день” — если число оканчивается на 1, кроме 11 → n % 10 == 1 && n % 100 != 11
// “дня” — если на 2, 3, 4, кроме 12, 13, 14 → (n % 10 == 2 || n % 10 == 3 || n % 10 == 4) && !(n % 100 >= 12 && n % 100 <= 14)
// “дней” — во всех остальных случаях (0, 5–20, 25–30 и т.д.)
func GetWateringInterval(wateringInterval int) string {
	// Проверяем последние цифры
	remainderFromDivisionByTen := wateringInterval % 10
	remainderFromDivisionByHundred := wateringInterval % 100

	switch {
	case remainderFromDivisionByTen == 1 && remainderFromDivisionByHundred != 11:
		return fmt.Sprintf("%d день", wateringInterval)
	case (remainderFromDivisionByTen == 2 || remainderFromDivisionByTen == 3 || remainderFromDivisionByTen == 4) &&
		(remainderFromDivisionByHundred < 12 || remainderFromDivisionByHundred > 14):
		return fmt.Sprintf("%d дня", wateringInterval)
	default:
		return fmt.Sprintf("%d дней", wateringInterval)
	}
}
