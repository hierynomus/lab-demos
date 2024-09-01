package domain

// Action is a struct that represents a point of interest
type Action struct {
	Name     string `json:"name" yaml:"name"`
	Duration int    `json:"duration" yaml:"duration"`
}

type DinoActions struct {
	DinoName string   `json:"name" yaml:"name"`
	Actions  []Action `json:"actions" yaml:"actions"`
}

type DinoNotFound struct {
	Name string
}

func (e *DinoNotFound) Error() string {
	return "Dino not found: " + e.Name
}
