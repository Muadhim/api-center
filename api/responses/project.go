package responses

type Project struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	AuthorID  uint   `json:"author_id"`
	Members   []User `json:"members"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
