package business

import "github.com/Toflex/directory_v2/ent"

type CreateBusinessRequest struct {
	UserID   string
	Name     string
	Category string
	Email    string
	// Website  string
}

type createBusiness struct {
	Name     string
	Category string
	Email    string
	// Website  string
}

type businessResult struct {
	ent.Business
}
