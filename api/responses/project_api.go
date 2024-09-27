package responses

import "time"

type ProjectApi struct {
	ID              uint      `json:"id"`
	Name            string    `json:"name"`
	FolderID        uint      `json:"folder_id"`
	Author          User      `json:"author"`
	UpdateBy        *User     `json:"update_by"`
	Method          string    `json:"method"`
	Path            string    `json:"path"`
	Request         string    `json:"request"`
	Response        string    `json:"response"`
	Description     string    `json:"description"`
	ExampleRequest  string    `json:"example_request"`
	ExampleResponse string    `json:"example_response"`
	ProjectID       uint      `json:"project_id"`
	CreatedAt       time.Time ` json:"created_at"`
	UpdatedAt       time.Time ` json:"updated_at"`
}
