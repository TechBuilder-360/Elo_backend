package rbac

type Permission string

const (
	CanAddMember       Permission = "member.add"
	CanViewMember      Permission = "member.view"
	CanUpdateMember    Permission = "member.update"
	CanDeleteMember    Permission = "member.delete"
	CanViewWallet      Permission = "wallet.view"
	CanViewTransaction Permission = "transaction.view"
)

var AllPermissions = map[Permission]string{
	CanAddMember:       "Permission to add a manager",
	CanViewMember:      "Permission to read member information",
	CanUpdateMember:    "Permission to update member information",
	CanDeleteMember:    "Permission to delete a member",
	CanViewWallet:      "Permission to view wallet information",
	CanViewTransaction: "Permission to view transaction information",
}

var RolePermissions = map[Role][]Permission{
	SuperAdminRole: {"*"},
	AdminRole:      {CanAddMember, CanViewMember, CanUpdateMember, CanDeleteMember, CanViewWallet, CanViewTransaction},
	ManagerRole:    {CanViewMember, CanViewWallet, CanViewTransaction},
}
