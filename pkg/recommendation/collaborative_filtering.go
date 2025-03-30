package recommendation

import (
	"log"
	"math"
	"math/rand"
	"product-recommendations-go/internal/models"
	"sort"
	"time"
)

// RecommendProducts генерує рекомендації на основі гібридного підходу
func RecommendProducts(userID uint, likes []*models.UserLike, orders []*models.Order, allProducts []*models.Product, limit int) ([]*models.Product, []float64) {
	// Ініціалізуємо рекомендації
	var recommendations []*models.Product
	var scores []float64

	// 1. Спочатку спробуємо колаборативну фільтрацію
	collaborativeRecs, collaborativeScores := getCollaborativeRecommendations(userID, likes, orders, allProducts, limit)
	recommendations = append(recommendations, collaborativeRecs...)
	scores = append(scores, collaborativeScores...)

	log.Println("Count of collab records: ", len(recommendations))

	// Якщо колаборативна фільтрація не дала результатів, використовуємо контентну фільтрацію
	if len(recommendations) == 0 {
		log.Printf("No collaborative recommendations found, trying content-based recommendations")
		contentRecs, contentScores := getContentBasedRecommendations(userID, likes, orders, allProducts, limit)
		recommendations = append(recommendations, contentRecs...)
		scores = append(scores, contentScores...)
	}

	log.Println("Count of content fill records: ", len(recommendations))

	// Якщо контентна фільтрація не дала результатів, використовуємо популярні товари
	if len(recommendations) == 0 {
		log.Printf("No content-based recommendations found, trying popularity-based recommendations")
		popularRecs, popularScores := getPopularityBasedRecommendations(userID, likes, orders, allProducts, limit)
		recommendations = append(recommendations, popularRecs...)
		scores = append(scores, popularScores...)
	}

	log.Println("Count of popular: ", len(recommendations))

	// Якщо все ще немає рекомендацій, використовуємо випадкові товари
	if len(recommendations) == 0 {
		log.Printf("No popularity-based recommendations found, using random recommendations")
		randomRecs, randomScores := getRandomRecommendations(userID, likes, orders, allProducts, limit)
		recommendations = append(recommendations, randomRecs...)
		scores = append(scores, randomScores...)
	}

	log.Println("Count of rand: ", len(recommendations))

	log.Printf("Final recommendations count: %d, scores count: %d", len(recommendations), len(scores))
	return recommendations, scores
}

// getCollaborativeRecommendations використовує колаборативну фільтрацію
func getCollaborativeRecommendations(userID uint, likes []*models.UserLike, orders []*models.Order, allProducts []*models.Product, limit int) ([]*models.Product, []float64) {
	var recommendations []*models.Product
	var scores []float64

	// Створення мапи продуктів, які користувач вже лайкнув або купив
	userProductMap := make(map[uint]bool)

	// Додавання лайкнутих продуктів
	for _, like := range likes {
		if like.UserID == userID {
			userProductMap[like.ProductID] = true
		}
	}

	// Додавання куплених продуктів
	for _, order := range orders {
		if order.UserID == userID {
			for _, item := range order.Items {
				userProductMap[item.ProductID] = true
			}
		}
	}

	// Знаходження подібних користувачів
	userSimilarity := make(map[uint]float64)

	for _, like := range likes {
		if like.UserID != userID && userProductMap[like.ProductID] {
			userSimilarity[like.UserID] += 1.0
		}
	}

	// Сортування користувачів за подібністю
	type UserSim struct {
		UserID     uint
		Similarity float64
	}

	var userSims []UserSim
	for uid, sim := range userSimilarity {
		userSims = append(userSims, UserSim{uid, sim})
	}

	sort.Slice(userSims, func(i, j int) bool {
		return userSims[i].Similarity > userSims[j].Similarity
	})

	// Рекомендації на основі подібності
	recommendationScores := make(map[uint]float64)

	// Обмежуємо кількість подібних користувачів
	maxUsers := 10
	if len(userSims) < maxUsers {
		maxUsers = len(userSims)
	}

	// Додаємо продукти, які лайкнули схожі користувачі
	for i := 0; i < maxUsers; i++ {
		if i >= len(userSims) {
			break
		}

		similarUserID := userSims[i].UserID
		similarityScore := userSims[i].Similarity

		for _, like := range likes {
			if like.UserID == similarUserID && !userProductMap[like.ProductID] {
				recommendationScores[like.ProductID] += similarityScore
			}
		}
	}

	// Сортування продуктів за рейтингом рекомендації
	type ProductScore struct {
		Product *models.Product
		Score   float64
	}

	var productScores []ProductScore
	for pid, score := range recommendationScores {
		for _, product := range allProducts {
			if product.ID == pid {
				productScores = append(productScores, ProductScore{product, score})
				break
			}
		}
	}

	sort.Slice(productScores, func(i, j int) bool {
		return productScores[i].Score > productScores[j].Score
	})

	// Додаємо рекомендації з колаборативної фільтрації
	for _, ps := range productScores {
		recommendations = append(recommendations, ps.Product)
		scores = append(scores, ps.Score)

		if len(recommendations) >= limit {
			break
		}
	}

	return recommendations, scores
}

// getContentBasedRecommendations використовує контентну фільтрацію на основі категорій
func getContentBasedRecommendations(userID uint, likes []*models.UserLike, orders []*models.Order, allProducts []*models.Product, limit int) ([]*models.Product, []float64) {
	var recommendations []*models.Product
	var scores []float64

	// Створення мапи продуктів, які користувач вже лайкнув або купив
	userProductMap := make(map[uint]bool)
	categoryPreferences := make(map[string]float64)

	// Збираємо ціни лайкнутих продуктів для розрахунку подібності
	var likedProductPrices []float64
	var sumLikedPrices float64

	// Додавання лайкнутих продуктів і визначення вподобань за категоріями
	for _, like := range likes {
		if like.UserID == userID {
			userProductMap[like.ProductID] = true

			// Знаходимо продукт, щоб отримати його категорію та ціну
			for _, product := range allProducts {
				if product.ID == like.ProductID {
					categoryPreferences[product.Category] += 1.0
					likedProductPrices = append(likedProductPrices, product.Price)
					sumLikedPrices += product.Price
					break
				}
			}
		}
	}

	// Додавання куплених продуктів і оновлення переваг категорій
	for _, order := range orders {
		if order.UserID == userID {
			for _, item := range order.Items {
				userProductMap[item.ProductID] = true

				// Знаходимо продукт, щоб отримати його категорію
				for _, product := range allProducts {
					if product.ID == item.ProductID {
						// Покупки мають більшу вагу, ніж лайки
						categoryPreferences[product.Category] += 2.0
						likedProductPrices = append(likedProductPrices, product.Price)
						sumLikedPrices += product.Price
						break
					}
				}
			}
		}
	}

	// Якщо немає переваг за категоріями, повертаємо пустий список
	if len(categoryPreferences) == 0 {
		return nil, nil
	}

	// Рахуємо глобальну популярність товарів
	productPopularity := make(map[uint]int)
	for _, like := range likes {
		productPopularity[like.ProductID]++
	}

	// Оцінюємо продукти на основі переваг категорій
	type ProductScore struct {
		Product *models.Product
		Score   float64
	}

	var productScores []ProductScore

	for _, product := range allProducts {
		if userProductMap[product.ID] {
			continue
		}

		// Базовий рейтинг на основі категорії
		score := categoryPreferences[product.Category]

		if score > 0 {
			// Детальне логування базового рейтингу
			log.Printf("Product %d (%s): Base category score: %.2f",
				product.ID, product.Name, score)

			// 1. Зменшуємо вплив новизни
			daysSinceCreation := time.Since(product.CreatedAt).Hours() / 24
			recencyBonus := 1.0 / (1.0 + daysSinceCreation/30) * 0.2 // Зменшено вагу
			score += recencyBonus

			log.Printf("  + Recency bonus: %.2f", recencyBonus)

			// 2. Враховуємо популярність товару (збільшена вага)
			popularity := float64(productPopularity[product.ID]) * 0.3
			score += popularity

			log.Printf("  + Popularity bonus: %.2f (%d likes)",
				popularity, productPopularity[product.ID])

			// 3. Додаємо схожість за ціною
			if len(likedProductPrices) > 0 {
				avgLikedPrice := sumLikedPrices / float64(len(likedProductPrices))
				priceDiff := math.Abs(product.Price - avgLikedPrice)
				priceSimilarityFactor := math.Max(0, 1.0-priceDiff/avgLikedPrice/2)
				score += priceSimilarityFactor * 0.3

				log.Printf("  + Price similarity bonus: %.2f (product: %.2f, avg liked: %.2f)",
					priceSimilarityFactor*0.3, product.Price, avgLikedPrice)
			}

			log.Printf("  = Final score: %.2f", score)

			productScores = append(productScores, ProductScore{product, score})
		}
	}

	// Сортуємо за рейтингом (більший - вищий)
	sort.Slice(productScores, func(i, j int) bool {
		return productScores[i].Score > productScores[j].Score
	})

	// Обмежуємо кількість рекомендацій
	maxRecommendations := limit
	if len(productScores) < maxRecommendations {
		maxRecommendations = len(productScores)
	}

	// Формуємо фінальний список рекомендацій
	for i := 0; i < maxRecommendations; i++ {
		recommendations = append(recommendations, productScores[i].Product)
		scores = append(scores, productScores[i].Score)
	}

	return recommendations, scores
}

// getPopularityBasedRecommendations використовує популярність товарів
func getPopularityBasedRecommendations(userID uint, likes []*models.UserLike, orders []*models.Order, allProducts []*models.Product, limit int) ([]*models.Product, []float64) {
	var recommendations []*models.Product
	var scores []float64

	// Створення мапи продуктів, які користувач вже лайкнув або купив
	userProductMap := make(map[uint]bool)

	// Додавання лайкнутих продуктів
	for _, like := range likes {
		if like.UserID == userID {
			userProductMap[like.ProductID] = true
		}
	}

	// Додавання куплених продуктів
	for _, order := range orders {
		if order.UserID == userID {
			for _, item := range order.Items {
				userProductMap[item.ProductID] = true
			}
		}
	}

	// Підрахунок популярності товарів
	popularity := make(map[uint]int)

	for _, like := range likes {
		popularity[like.ProductID]++
	}

	// Якщо немає лайків взагалі, повертаємо пустий список
	if len(popularity) == 0 {
		return nil, nil
	}

	type PopularProduct struct {
		Product *models.Product
		Count   int
	}

	var popularProducts []PopularProduct
	for pid, count := range popularity {
		for _, product := range allProducts {
			if product.ID == pid {
				popularProducts = append(popularProducts, PopularProduct{product, count})
				break
			}
		}
	}

	// Сортуємо за популярністю
	sort.Slice(popularProducts, func(i, j int) bool {
		return popularProducts[i].Count > popularProducts[j].Count
	})

	// Додаємо популярні товари до рекомендацій (окрім тих, що вже лайкав користувач)
	for _, pp := range popularProducts {
		if !userProductMap[pp.Product.ID] {
			recommendations = append(recommendations, pp.Product)
			scores = append(scores, float64(pp.Count))
		}

		// Обмежуємо кількість рекомендацій
		if len(recommendations) >= limit {
			break
		}
	}

	return recommendations, scores
}

// getRandomRecommendations генерує випадкові рекомендації
func getRandomRecommendations(userID uint, likes []*models.UserLike, orders []*models.Order, allProducts []*models.Product, count int) ([]*models.Product, []float64) {
	var recommendations []*models.Product
	var scores []float64

	// Створення мапи продуктів, які користувач вже лайкнув або купив
	userProductMap := make(map[uint]bool)

	// Додавання лайкнутих продуктів
	for _, like := range likes {
		if like.UserID == userID {
			userProductMap[like.ProductID] = true
		}
	}

	// Додавання куплених продуктів
	for _, order := range orders {
		if order.UserID == userID {
			for _, item := range order.Items {
				userProductMap[item.ProductID] = true
			}
		}
	}

	// Ініціалізуємо випадковий генератор
	rand.NewSource(time.Now().UnixNano())

	// Перемішуємо всі товари
	shuffledProducts := make([]*models.Product, len(allProducts))
	copy(shuffledProducts, allProducts)
	rand.Shuffle(len(shuffledProducts), func(i, j int) {
		shuffledProducts[i], shuffledProducts[j] = shuffledProducts[j], shuffledProducts[i]
	})

	// Додаємо перші count (або менше) товарів
	maxRandomProducts := count
	if len(shuffledProducts) < maxRandomProducts {
		maxRandomProducts = len(shuffledProducts)
	}

	for i := 0; i < len(shuffledProducts) && len(recommendations) < maxRandomProducts; i++ {
		// Пропускаємо товари, які користувач уже лайкав/купував
		if !userProductMap[shuffledProducts[i].ID] {
			recommendations = append(recommendations, shuffledProducts[i])
			// Додаємо фіксований низький рейтинг для випадкових рекомендацій
			scores = append(scores, 0.1)
		}
	}

	return recommendations, scores
}
