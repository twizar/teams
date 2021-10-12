package entity

type Team struct {
	id     string
	name   string
	league string
	rating float64
}

func (t Team) ID() string {
	return t.id
}

func (t Team) Name() string {
	return t.name
}

func (t Team) League() string {
	return t.league
}

func (t Team) Rating() float64 {
	return t.rating
}

func NewTeam(id, name, league string, rating float64) *Team {
	return &Team{id: id, name: name, league: league, rating: rating}
}
