package entity

import (
	"time"
)


type Jobs struct{
	Id string
	Client_id string
	Name string
	Comp_name string
	Status bool
	StartDate string
	Location string
	EndDate string
	Created_at time.Time
	Updated_at time.Time
	Deleted_at time.Time
}
