package entity

type Flat struct {
	ID      int    `json:"id"`
	Number  int    `json:"number"`
	HouseID int    `json:"house_id"`
	Price   int    `json:"price"`
	Rooms   int    `json:"rooms"`
	Status  string `json:"status"` // можно сделать через iota, чтобы задать конкретные варианты
}
