package service

import (
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewRoleService, NewOrganizationService, NewUserService, NewMenuService,
	NewEmployeeService, NewBankService)
