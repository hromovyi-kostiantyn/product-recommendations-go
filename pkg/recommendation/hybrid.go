// Package recommendation реалізує алгоритми рекомендаційної системи
// для генерації персоналізованих рекомендацій товарів на основі
// поведінки користувачів та характеристик товарів.
//
// Пакет містить реалізації різних підходів до формування рекомендацій:
//   - Колаборативна фільтрація (User-based Collaborative Filtering)
//   - Фільтрація на основі вмісту (Content-based Filtering)
//   - Гібридні алгоритми, що поєднують різні підходи
//
// Приклад використання колаборативної фільтрації:
//
//	userRatings := map[int64]map[int64]float64{
//		1: {101: 5.0, 102: 3.0, 103: 2.5},
//		2: {101: 4.0, 103: 4.5, 104: 3.5},
//		3: {102: 4.5, 104: 4.0, 105: 3.5},
//	}
//	prediction := recommendation.PredictRating(1, 104, userRatings)
package recommendation

import "math"

// SimilarityMetric визначає інтерфейс для обчислення подібності між елементами.
// Реалізації цього інтерфейсу використовуються для порівняння користувачів
// або продуктів у рекомендаційних алгоритмах.
type SimilarityMetric interface {
	// Calculate обчислює подібність між двома векторами.
	// Повертає значення у діапазоні від 0 до 1, де 1 означає
	// повну ідентичність, а 0 - повну відмінність.
	Calculate(vector1, vector2 []float64) float64
}

// CosineSimilarity реалізує метрику подібності на основі косинуса кута
// між двома векторами. Це стандартний метод вимірювання подібності
// у рекомендаційних системах з колаборативною фільтрацією.
type CosineSimilarity struct{}

// Calculate обчислює косинусну подібність між двома векторами.
// Повертає значення у діапазоні від 0 до 1.
//
// Косинусна подібність обчислюється за формулою:
//
//	similarity = (A · B) / (||A|| × ||B||)
//
// де A · B - скалярний добуток векторів,
// ||A|| і ||B|| - евклідові норми (довжини) векторів.
//
// Параметри:
//   - vector1: перший вектор числових значень
//   - vector2: другий вектор числових значень (повинен мати таку ж довжину як vector1)
//
// Повертає:
//   - подібність як число від 0 до 1
//
// Якщо вектори мають різну довжину або один з векторів має нульову довжину,
// функція повертає 0.
func (cs CosineSimilarity) Calculate(vector1, vector2 []float64) float64 {
	if len(vector1) != len(vector2) || len(vector1) == 0 {
		return 0
	}

	var dotProduct, magnitude1, magnitude2 float64

	for i := 0; i < len(vector1); i++ {
		dotProduct += vector1[i] * vector2[i]
		magnitude1 += vector1[i] * vector1[i]
		magnitude2 += vector2[i] * vector2[i]
	}

	magnitude1 = math.Sqrt(magnitude1)
	magnitude2 = math.Sqrt(magnitude2)

	if magnitude1 == 0 || magnitude2 == 0 {
		return 0
	}

	return dotProduct / (magnitude1 * magnitude2)
}

// PredictRating прогнозує рейтинг, який користувач може дати певному товару
// на основі оцінок інших користувачів використовуючи колаборативну фільтрацію.
//
// Алгоритм працює в кілька етапів:
//  1. Знаходить користувачів, які оцінили цільовий товар
//  2. Обчислює подібність між цільовим користувачем і кожним з цих користувачів
//  3. Обчислює зважену суму оцінок товару, використовуючи подібність як ваги
//
// Параметри:
//   - targetUser: ідентифікатор користувача, для якого прогнозується рейтинг
//   - targetItem: ідентифікатор товару, для якого прогнозується рейтинг
//   - userRatings: двовимірна карта, що містить оцінки товарів користувачами
//     (перший ключ - ID користувача, другий ключ - ID товару, значення - оцінка)
//
// Повертає:
//   - прогнозований рейтинг товару для користувача як число з плаваючою точкою
//
// Якщо неможливо зробити прогноз (наприклад, відсутні подібні користувачі),
// функція повертає 0.
func PredictRating(targetUser int64, targetItem int64, userRatings map[int64]map[int64]float64) float64 {
	// Отримуємо оцінки цільового користувача
	targetUserRatings := userRatings[targetUser]

	// Перевіряємо, чи не оцінив користувач уже цей товар
	if _, exists := targetUserRatings[targetItem]; exists {
		return targetUserRatings[targetItem]
	}

	// Знаходимо користувачів, які оцінили цільовий товар
	similars := make(map[int64]float64)
	for userID, ratings := range userRatings {
		if userID == targetUser {
			continue
		}
		if _, hasRated := ratings[targetItem]; hasRated {
			// Обчислюємо подібність між користувачами
			similarity := calculateUserSimilarity(targetUser, userID, userRatings)
			if similarity > 0 {
				similars[userID] = similarity
			}
		}
	}

	// Якщо немає подібних користувачів, неможливо зробити прогноз
	if len(similars) == 0 {
		return 0
	}

	// Обчислюємо зважену суму оцінок
	var weightedSum, similaritySum float64
	for userID, similarity := range similars {
		weightedSum += similarity * userRatings[userID][targetItem]
		similaritySum += similarity
	}

	// Повертаємо прогнозований рейтинг
	if similaritySum == 0 {
		return 0
	}
	return weightedSum / similaritySum
}

// calculateUserSimilarity обчислює подібність між двома користувачами
// на основі їхніх оцінок спільних товарів.
//
// Параметри:
//   - user1: ідентифікатор першого користувача
//   - user2: ідентифікатор другого користувача
//   - userRatings: двовимірна карта, що містить оцінки товарів користувачами
//
// Повертає:
//   - подібність як число від 0 до 1
func calculateUserSimilarity(user1, user2 int64, userRatings map[int64]map[int64]float64) float64 {
	ratings1 := userRatings[user1]
	ratings2 := userRatings[user2]

	// Знаходимо спільні товари
	commonItems := make(map[int64]bool)
	for itemID := range ratings1 {
		if _, ok := ratings2[itemID]; ok {
			commonItems[itemID] = true
		}
	}

	// Якщо немає спільних товарів, користувачі не подібні
	if len(commonItems) == 0 {
		return 0
	}

	// Створюємо вектори оцінок для спільних товарів
	vector1 := make([]float64, 0, len(commonItems))
	vector2 := make([]float64, 0, len(commonItems))

	for itemID := range commonItems {
		vector1 = append(vector1, ratings1[itemID])
		vector2 = append(vector2, ratings2[itemID])
	}

	// Обчислюємо косинусну подібність
	similarity := CosineSimilarity{}
	return similarity.Calculate(vector1, vector2)
}
