package adaptivecard

const (
	TypeTextBlock      = "TextBlock"
	TypeColumnSet      = "ColumnSet"
	TypeColumn         = "Column"
	TypeActionSet      = "ActionSet"
	TypeActionShowCard = "Action.ShowCard"
	TypeActionOpenUrl  = "Action.OpenUrl"
	TypeAdaptiveCard   = "AdaptiveCard"

	WeightBorder  = "Border"
	WeightLighter = "Lighter"
	WeightDefault = "Default"

	WHAutomatic = "auto"
	WHStretch   = "stretch"

	SizeLarge  = "Large"
	SizeMedium = "Medium"

	ColorDefault   = "Default"
	ColorDark      = "Dark"
	ColorLight     = "Light"
	ColorAccent    = "Accent"
	ColorGood      = "Good"
	ColorWarning   = "Warning"
	ColorAttention = "Attention"
)

// type Attachment struct {
// 	ContentType string              `json:"contentType"`
// 	ContentUrl  interface{}         `json:"contentUrl"`
// 	Content     AdaptiveCardContent `json:"content"`
// }

type AdaptiveCardContent struct {
	Type                     string             `json:"type"`
	Schema                   string             `json:"$schema"`
	Version                  string             `json:"version"`
	Body                     []AdaptiveCardBody `json:"body"`
	VerticalContentAlignment string             `json:"verticalContentAlignment"`
	MSTeams                  MSTeams            `json:"msteams"`
}

// type MSTeams struct {
// 	Entities []MSTeamsEntity `json:"entities"`
// }
//
// type MSTeamsEntity struct {
// 	Type      string                 `json:"type"`
// 	Text      string                 `json:"text"`
// 	Mentioned MSTeamsEntityMentioned `json:"mentioned"`
// }
//
// type MSTeamsEntityMentioned struct {
// 	ID   string `json:"id"`
// 	Name string `json:"name"`
// }

type AdaptiveCardBody struct {
	Type    string   `json:"type"`
	Size    string   `json:"size,omitempty"`
	Weight  string   `json:"weight,omitempty"`
	Text    string   `json:"text,omitempty"`
	Color   string   `json:"color,omitempty"`
	Wrap    bool     `json:"wrap,omitempty"`
	Columns []Column `json:"columns,omitempty"`
}

type Column struct {
	Type  string      `json:"type"`
	Width interface{} `json:"width"`
	Items []Item      `json:"items"`
}

type Item struct {
	Type    string   `json:"type"`
	Text    string   `json:"text,omitempty"`
	Wrap    bool     `json:"wrap,omitempty"`
	Weight  string   `json:"weight,omitempty"`
	Color   string   `json:"color,omitempty"`
	Actions []Action `json:"actions,omitempty"`
}

type Action struct {
	Type  string     `json:"type"`
	Title string     `json:"title"`
	Url   string     `json:"url,omitempty"`
	Card  ActionCard `json:"card,omitempty"`
}

type ActionCard struct {
	Type string           `json:"type"`
	Body []ActionCardBody `json:"body"`
}

type ActionCardBody struct {
	Type string `json:"type"`
	Text string `json:"text"`
	Wrap bool   `json:"wrap"`
}

type Condition struct {
	Value  []string
	Volume string
}
