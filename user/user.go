package user

type User struct {
	ID     string `sql:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Credit int
}
