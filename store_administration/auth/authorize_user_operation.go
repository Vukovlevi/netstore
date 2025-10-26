package auth

import "github.com/vukovlevi/netstore/store_administration/model"

func CanUserSetRole(user model.User, role int) bool {
    if user.Role == ROLE_HR && role == ROLE_STORE_LEADER_ID {
        return false
    }
    return true
}

func CanUserDisablePasswordChange(user model.User) bool {
    if user.Role == ROLE_HR {
        return false
    }
    return true
}
