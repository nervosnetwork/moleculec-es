package generator

import (
	"fmt"
)

type Options struct {
	HasBigInt bool
}

type Field struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type Declaration struct {
	Type      string       `json:"type"`
	Name      string       `json:"name"`
	Item      string       `json:"item"`
	Items     []UnionItems `json:"items"`
	ItemCount int          `json:"item_count"`
	Fields    []Field      `json:"fields"`
}

type UnionItems struct {
	Type string `json:"typ"`
	Id   uint64 `json:"id"`
}

type SyntaxVersion struct {
	Version uint64 `json:"version"`
}

type Schema struct {
	SyntaxVersion SyntaxVersion `json:"syntax_version"`
	Namespace     string        `json:"namespace"`
	Declarations  []Declaration `json:"declarations"`
}

func (s Schema) FindDeclaration(itemName string) (Declaration, error) {
	for _, declaration := range s.Declarations {
		if declaration.Name == itemName {
			return declaration, nil
		}
	}
	return Declaration{}, fmt.Errorf("Cannot find type %s!", itemName)
}
