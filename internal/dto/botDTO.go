package dto

type Messages struct {
	Role    string   `json:"role"`
	Content string   `json:"content"`
	Images  []string `json:"images"`
}

type ImageRequest struct {
	Model    string     `json:"model"`
	Messages []Messages `json:"messages"`
	Stream   bool       `json:"stream"`
}
