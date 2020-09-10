package parser

// PokiPokiProperty represents a type definition of an object's property
type PokiPokiProperty struct {
	Name string
	Type []string
}

// PokiPokiObject represents a type definition of an object
type PokiPokiObject struct {
	Name       string
	Properties []PokiPokiProperty
	Children   []string
	Edges      map[string]PokiPokiRelationship
}

// PokiPokiDocument represents the parsed form of a pokipoki file
type PokiPokiDocument struct {
	SchemaName    string
	SchemaVersion int64
	Objects       map[string]PokiPokiObject
}

type RelationshipDirection int

const (
	ToItem RelationshipDirection = iota
	FromItem
)

// PokiPokiRelationship represents a relationship
type PokiPokiRelationship struct {
	Name      string
	Direction RelationshipDirection
	Unique    bool
	From      struct {
		Type string
		Name string
	}
	To struct {
		Type string
	}
}
