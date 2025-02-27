package blocks

type BlockRunner interface {
  Run()
}

type Block struct {
  BlockRunner
  FullText string `json:"full_text"`
  ShortText string  `json:"short_text,omitempty"`
  Color string  `json:"color,omitempty"`
  Background string `json:"background,omitempty"`
  Border string `json:"border,omitempty"`
  BorderTop int `json:"border_top,omitempty"`
  BorderBottom int `json:"border_bottom,omitempty"`
  BorderLeft int `json:"border_left,omitempty"`
  BorderRight int `json:"border_right,omitempty"`
  MinWidth int `json:"min_width,omitempty"`
  Align string `json:"align,omitempty"`
  Name string `json:"name,omitempty"`
  Instance string `json:"instance,omitempty"`
  Urgent bool `json:"urgent,omitempty"`
  Separator bool `json:"separator,omitempty"`
  SeparatorBlockWidth int `json:"separator_block_width,omitempty"`
  Markup string `json:"markup,omitempty"`
}
