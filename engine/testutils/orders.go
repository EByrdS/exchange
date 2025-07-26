package testutils

import (
	"exchange/engine/order"
	"fmt"
	"math/rand/v2"
)

func OrdersRandom(pair string, side order.OrderSide, minPrice, maxPrice, minVolume, maxVolume, count uint64) []*order.Order {
	orders := make([]*order.Order, 0, count)
	for i := range count {
		price := rand.Uint64N(maxPrice-minPrice) + minPrice
		volume := rand.Uint64N(maxVolume-minVolume) + minVolume
		orders = append(orders, &order.Order{
			Pair:   pair,
			ID:     fmt.Sprintf("%d", i),
			Side:   side,
			Price:  price,
			Volume: volume,
		})
	}

	return orders
}

func OrdersDeterministic(pair string, side order.OrderSide, minPrice, maxPrice, minVolume, maxVolume, count uint64) []*order.Order {
	orders := make([]*order.Order, 0, count)
	price := minPrice
	volume := minVolume
	for i := range count {
		orders = append(orders, &order.Order{
			Pair:   pair,
			ID:     fmt.Sprintf("%d", i),
			Side:   side,
			Price:  price,
			Volume: volume,
		})

		price += 1
		if price > maxPrice {
			price = minPrice
		}

		volume += 1
		if volume > maxVolume {
			volume = minVolume
		}
	}

	return orders
}
