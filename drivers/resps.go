package drivers

type power struct {
	Power string `json:"power"`
}

type blanked struct {
	Blanked bool `json:"blanked"`
}

type input struct {
	Input string `json:"input"`
}

type muted struct {
	Muted bool `json:"muted"`
}

type volume struct {
	Volume int `json:"volume"`
}

type activeSignal struct {
	ActiveSignal bool `json:"activeSignal"`
}
