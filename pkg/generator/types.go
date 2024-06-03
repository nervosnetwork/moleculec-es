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

type DeclarationOld struct {
	Type      string   `json:"type"`
	Name      string   `json:"name"`
	Item      string   `json:"item"`
	Items     []string `json:"items"`
	ItemCount int      `json:"item_count"`
	Fields    []Field  `json:"fields"`
}

type SyntaxVersion struct {
	Version uint64 `json:"version"`
}

type Schema struct {
	SyntaxVersion SyntaxVersion `json:"syntax_version"`
	Namespace     string        `json:"namespace"`
	Declarations  []Declaration `json:"declarations"`
}

type SchemaOld struct {
	Namespace    string           `json:"namespace"`
	Declarations []DeclarationOld `json:"declarations"`
}

func (s Schema) FindDeclaration(itemName string) (Declaration, error) {
	for _, declaration := range s.Declarations {
		if declaration.Name == itemName {
			return declaration, nil
		}
	}
	return Declaration{}, fmt.Errorf("Cannot find type %s!", itemName)
}

func (s SchemaOld) ChangeToNew() Schema {
	var de []Declaration
	for _, item := range s.Declarations {
		de = append(de, item.ChangeToNew())
	}
	return Schema{
		SyntaxVersion: SyntaxVersion{Version: 1},
		Namespace:     s.Namespace,
		Declarations:  de,
	}
}

func (de DeclarationOld) ChangeToNew() Declaration {
	var unionItems []UnionItems

	for i, item := range de.Items {
		unionItems = append(unionItems, UnionItems{Id: uint64(i), Type: item})
	}

	return Declaration{Type: de.Type, Name: de.Name, Item: de.Item, Items: unionItems, Fields: de.Fields}
}
