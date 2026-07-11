package rbac

type Role string

const (
	SuperAdminRole Role = "Super Admin"
	AdminRole      Role = "Admin"
	ManagerRole    Role = "Manager"
)

func (r Role) ToString() string {
	return string(r)
}
