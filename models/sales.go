package models

type SalesSummary struct {
	TotalSales        int `json:"total_sales"`
	TotalTransactions int `json:"total_transactions"`
}

type BestSellingProduct struct {
	ProductID   int    `json:"product_id"`
	ProductName string `json:"product_name"`
	TotalSold   int    `json:"total_sold"`
}

type SalesReport struct {
	Summary            SalesSummary         `json:"summary"`
	BestSellingProducts []BestSellingProduct `json:"best_selling_products"`
}
