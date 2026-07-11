package manager

type Manager struct {
	UserID     string `json:"user_id"`
	BusinessID string `json:"business_id"`
	IsOwner    bool   `json:"is_owner"`
	RoleID     string `json:"role_id"`
}
