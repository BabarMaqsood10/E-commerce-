package cart

import (
	"fmt"

	"myproject/types"
)

// getProductIDs → helper function that takes a slice of CartItemPayload structs and returns a slice of product IDs extracted from the items
func getCartItemIDs(items []types.CartItemPayload) ([]int, error) {
	// Create a slice to hold the product IDs, with the same length as the items slice
	productIDs := make([]int, len(items))
	// Iterate through the items slice, extracting the ProductID from each CartItemPayload and storing it in the productIDs slice. Return an error if any item has a non-positive quantity
	for i, item := range items {
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity for product ID %d: %d", item.ProductID, item.Quantity)
		}
		// Store the ProductID in the productIDs slice
		productIDs[i] = item.ProductID
	}
	// Return the slice of product IDs and a nil error, indicating successful extraction of product IDs from the cart items
	return productIDs, nil
}

func (h *Handler) createOrder(ps []types.Product, items []types.CartItemPayload, userID int) (int, float64, error) {
	// Create a map to hold the products by their ID for easy lookup
	// loop through the products and populate the map with product ID as the key and the Product struct as the value
	productMap := make(map[int]types.Product)
	for _, p := range ps {
		productMap[p.ID] = p
	}

	// check if the product is available in the store
	// check the products quantity in the store. If the quantity is less than the requested quantity, return an error indicating that the product is out of stock
	if err := checkIfCartIsInStock(items, productMap); err != nil {
		return 0, 0, err
	}

	// calculate the total price
	totalPrice := calculateTotalPrice(items, productMap)

	// reduce the quantity of the product in the store
	// loop through the items in the cart, reduce the quantity of each product in the productMap by the quantity requested in the cart, and update the product in the store using the UpdateProduct method. This ensures that the stock levels are accurately reflected after the order is created
	for _, item := range items {
		product := productMap[item.ProductID]
		product.Quantity -= item.Quantity
		productMap[item.ProductID] = product
		h.productStore.UpdateProduct(product)
	}

	// create the order
	orderID, err := h.store.CreateOrder(types.Order{
		UserID:  userID,
		Total:   totalPrice,
		Status:  "pending",
		Address: "123 Main St",
	})
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create order: %v", err)
	}

	// create the order items
	for _, item := range items {
		h.store.CreateOrderItem(types.OrderItem{
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     productMap[item.ProductID].Price,
		})
	}

	return orderID, totalPrice, nil
}

// checkIfCartIsInStock → checks if the items in the cart are available in stock, returns an error if any item is out of stock or if the cart is empty
func checkIfCartIsInStock(cartItems []types.CartItemPayload, products map[int]types.Product) error {
	// Check if the cart is empty, returning an error if there are no items in the cart
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}
	// Iterate through the cart items, checking if each product exists in the products map and if the requested quantity is available. Return an error if any product does not exist or if the requested quantity exceeds the available stock
	for _, item := range cartItems {
		// Check if the product exists in the products map
		product, exists := products[item.ProductID]
		if !exists {
			return fmt.Errorf("product with ID %d does not exist", item.ProductID)
		}
		// Check if the requested quantity is available in stock
		if product.Quantity < item.Quantity {
			return fmt.Errorf("product with ID %d is out of stock", product.ID)
		}
	}

	return nil
}

// calculateTotalPrice → calculates the total price of the items in the cart based on their quantities and prices, returning the total price as a float64
func calculateTotalPrice(cartItems []types.CartItemPayload, products map[int]types.Product) float64 {
	var totalPrice float64
	// Iterate through the cart items, calculating the total price by multiplying the quantity of each item by its corresponding product price and accumulating the result
	for _, item := range cartItems {
		product := products[item.ProductID]
		totalPrice += float64(item.Quantity) * product.Price
	}
	// Return the calculated total price for the items in the cart
	return totalPrice
}
