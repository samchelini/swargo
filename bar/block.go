package bar

import (
  "encoding/json"
)

type BlockRunner interface {
	Run()
	Sync(updater chan bool)
	Update()
  String() string
}

type Block struct {
	BlockRunner `json:"-"`
	FullText            string `json:"full_text"`
	ShortText           string `json:"short_text,omitempty"`
	Color               string `json:"color,omitempty"`
	Background          string `json:"background,omitempty"`
	Border              string `json:"border,omitempty"`
	BorderTop           int    `json:"border_top,omitempty"`
	BorderBottom        int    `json:"border_bottom,omitempty"`
	BorderLeft          int    `json:"border_left,omitempty"`
	BorderRight         int    `json:"border_right,omitempty"`
	MinWidth            int    `json:"min_width,omitempty"`
	Align               string `json:"align,omitempty"`
	Name                string `json:"name,omitempty"`
	Instance            string `json:"instance,omitempty"`
	Urgent              bool   `json:"urgent,omitempty"`
	Separator           bool   `json:"separator,omitempty"`
	SeparatorBlockWidth int    `json:"separator_block_width,omitempty"`
	Markup              string `json:"markup,omitempty"`
	updater             chan bool
}

func (b *Block) Sync(updater chan bool) {
	b.updater = updater
}

func (b *Block) Update() {
	b.updater <- true
}

func (b *Block) String() string {
    json, _ := json.Marshal(b)
    return string(json)
}
