package entities

type Branch struct {
	Id   uint64 `bson:"id" json:"id"`
	Name string `bson:"name" json:"name"`
}

func (Branch) CollectionName() string {
	return "branches"
}

type Department struct {
	Id   uint64 `bson:"id" json:"id"`
	Name string `bson:"name" json:"name"`
}

func (Department) CollectionName() string {
	return "departments"
}
