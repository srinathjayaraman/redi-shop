package stock

type Stock struct {
	ID     string `sql:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Price  int
	Number int
}
