package dto

type ProductResponse struct {
	Name           string   `json:"name"`                     // @NotBlank + pattern
	Category       string   `json:"category"`                 // @NotBlank + pattern
	Brand          string   `json:"brand"`                    // @NotBlank + pattern
	ExpirationDate *string  `json:"expirationDate,omitempty"` // LocalDate yyyy-MM-dd
	Quantity       *string  `json:"quantity,omitempty"`       // @Min(0)
	Unit           *string  `json:"unit,omitempty"`
	Location       *string  `json:"location,omitempty"`
	Barcode        *string  `json:"barcode,omitempty"`
	Notes          *string  `json:"notes,omitempty"`
	Status         Status   `json:"status,omitempty"`   // Enum
	Priority       Priority `json:"priority,omitempty"` // Enum
}

type Status string

const (
	StatusNew    Status = "NEW"
	StatusOpened Status = "OPENED"
	StatusUsed   Status = "USED"
	StatusEmpty  Status = "EMPTY"
)

type Priority string

const (
	PriorityHigh   Priority = "HIGH"
	PriorityMedium Priority = "MEDIUM"
	PriorityLow    Priority = "LOW"
)
