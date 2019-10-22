package drivers

type Power struct {
	Power string `json:"power"`
}

type Blanked struct {
	Blanked bool `json:"blanked"`
}

type Input struct {
	Input string `json:"input"`
}

type Muted struct {
	Muted bool `json:"muted"`
}

type Volume struct {
	Volume int `json:"volume"`
}

type ActiveSignal struct {
	ActiveSignal bool `json:"activeSignal"`
}
