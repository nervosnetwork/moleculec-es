package generator

type Field struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type Declaration struct {
	Type      string   `json:"type"`
	Name      string   `json:"name"`
	Item      string   `json:"item"`
	Items     []string `json:"items"`
	ItemCount int      `json:"item_count"`
	Fields    []Field  `json:"fields"`
}

type Schema struct {
	Namespace    string        `json:"namespace"`
	Declarations []Declaration `json:"declarations"`
}
