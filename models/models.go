package models

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Job struct {
	ID           int    `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Requirements string `json:"requirements"`
	EmployerID   int    `json:"employer_id"`
}

type Application struct {
	ID       int    `json:"id"`
	JobID    int    `json:"job_id"`
	TalentID int    `json:"talent_id"`
	Status   string `json:"status"`
}
