package models

type Phase string

const (
	PhaseExpansion   Phase = "Expansion"
	PhaseContraction Phase = "Contraction"
	PhaseTransition  Phase = "Transition"
)

type PhaseResponse struct {
	Asset      string  `json:"asset"`
	Phase      Phase   `json:"phase"`
	Confidence float64 `json:"confidence"`
	Status     Status  `json:"status"`
}
