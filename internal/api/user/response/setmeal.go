package response

type DishItemVO struct {
	Copies      int    `json:"copies"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Name        string `json:"name"`
}
