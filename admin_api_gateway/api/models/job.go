package models

type Job struct {
	Id string
	Client_id string
	Name string
	CompName string
	Status bool
	StartDate string
	EndDate string
	Location string
	Created_at string
	Updated_at string
}

type JobCreateReq struct {
	Client_id string
	Name string
	CompName string
	Status bool
	StartDate string
	EndDate string
	Location string
}

type JobList struct{
	ListJob []Job
}
type JobUpdateRequest struct{
	Id string
	Client_id string
	Name string
	CompName string
	Status bool
	StartDate string
	EndDate string
	Location string
}
type Owner struct {
	Id                   string   
	Name                 string   
	LastName             string
	Email                string
	Password             string   
	CreatedAt            string   
	UpdatedAt            string   
	DeletedAt            string   
	RefreshToken         string   
	Role                 string   
	
}
type JobWithOwner struct {
	Id                   string   `json:"Id"`
	ClientId             string   `json:"Client_id"`
	Name                 string   `json:"Name"`
	CompName             string   `json:"Comp_name"`
	Status               bool     `json:"Status"`
	StartDate            string   `json:"StartDate"`
	Location             string   `json:"Location"`
	EndDate              string   `json:"EndDate"`
	CreatedAt            string   `json:"Created_at"`
	UpdatedAt            string   `json:"Updated_at"`
	Owner                *Owner   `json:"owner"`
	
}


type GetAllJobByClientIdResponse struct {
	Jobs                 []*JobWithOwner `json:"jobs"`
}