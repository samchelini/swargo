package bar

// methods all blocks should implement
// only Run() needs to be implemented if using the BlockTemplate
type Block interface {
	Run()
	Sync(update chan bool, err chan string)
	Update()
}

// contains the block fields defined by the swaybar protocol
// implements the Update, Sync, and String functions for the Block interface
type BlockTemplate struct {
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
	update              chan bool
	err                 chan string
}

// adds the bar's update channel to the block
func (b *BlockTemplate) Sync(update chan bool, err chan string) {
	b.update = update
	b.err = err
}

// sends update signal to the update channel
// this triggers the bar to update the status line
func (b *BlockTemplate) Update() {
	b.update <- true
}
