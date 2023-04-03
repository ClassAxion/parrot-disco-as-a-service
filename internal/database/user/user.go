package user

type Location struct {
	Latitude  float64
	Longitude float64
	Altitude  int
}

type User struct {
	ID                int
	Email             string
	Password          string
	Name              string
	Hash              *string
	ZeroTierNetworkId *string
	ZeroTierDiscoIP   *string
	HomeLocation      *Location
}
