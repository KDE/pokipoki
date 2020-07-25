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
}

// PokiPokiDocument represents the parsed form of a pokipoki file
type PokiPokiDocument struct {
	Objects map[string]PokiPokiObject
}
