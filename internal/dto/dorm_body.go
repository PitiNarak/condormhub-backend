package dto

type DormRequestBody struct {
	Name      string  `json:"name" validate:"required"`
	Size      float64 `json:"size" validate:"required,gt=0"`
	Bedrooms  int     `json:"bedrooms" validate:"required,gte=0"`
	Bathrooms int     `json:"bathrooms" validate:"required,gte=0"`
	Address   struct {
		District    string `json:"district" validate:"required"`
		Subdistrict string `json:"subdistrict" validate:"required"`
		Province    string `json:"province" validate:"required"`
		Zipcode     string `json:"zipcode" validate:"required,numeric,len=5"`
	} `json:"address" validate:"required"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Description string  `json:"description"`
}
