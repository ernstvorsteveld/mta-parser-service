package product

type Product struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	Description string `json:"description"`
}

func GetProducts() []Product {
	return []Product{
		Product{100, "BassTune Headset 2.0", 200, "A headphone with a inbuilt high-quality microphone"},
		Product{101, "Fastlane Toy Car", 100, "A toy car that comes with a free HD camera"},
		Product{101, "ATV Gear Mouse", 75, "A high-quality mouse for office work and gaming"},
	}
}
