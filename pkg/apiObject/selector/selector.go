package selector

type Selector struct{
	MatchLabels map[string]string `json:"matchLabels" yaml:"matchLabels"`
}