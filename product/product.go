package product

type Product struct {
	Id int `json:"id"`
	Name string `json:"name"`
	Price int `json:"price"`
	Description string `json:"description"`
}