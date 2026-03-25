package dto

type ProductInfo struct {
	Product struct {
		Brands              string `json:"brands"`
		ExpirationDate      string `json:"expiration_date"`
		ProductName         string `json:"product_name"`
		ProductQuantityUnit string `json:"product_quantity_unit"`
		ProductType         string `json:"product_type"`
		Quantity            string `json:"quantity"`
	} `json:"product"`
}
