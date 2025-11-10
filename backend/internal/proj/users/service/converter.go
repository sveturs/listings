// backend/internal/proj/users/service/converter.go
package service

import (
	"backend/internal/domain/models"

	"github.com/sveturs/auth/pkg/entity"
)

// ToEntityUpdateProfile converts models.UserProfileUpdate to entity.UpdateProfileRequest
func ToEntityUpdateProfile(update *models.UserProfileUpdate) entity.UpdateProfileRequest {
	req := entity.UpdateProfileRequest{}

	if update.Name != nil {
		req.Name = update.Name
	}
	if update.Phone != nil {
		req.Phone = update.Phone
	}
	if update.Bio != nil {
		req.Bio = update.Bio
	}
	if update.Timezone != nil {
		req.Timezone = update.Timezone
	}
	if update.City != nil {
		req.City = update.City
	}
	if update.Country != nil {
		req.Country = update.Country
	}

	return req
}

// FromEntityUserProfile converts entity.UserProfile to models.User
func FromEntityUserProfile(profile *entity.UserProfile) *models.User {
	if profile == nil {
		return nil
	}

	return &models.User{
		ID:    profile.ID,
		Email: profile.Email,
		Name:  profile.Name,
	}
}

// FromEntityUserProfileToProfile converts entity.UserProfile to models.UserProfile
func FromEntityUserProfileToProfile(profile *entity.UserProfile) *models.UserProfile {
	if profile == nil {
		return nil
	}

	var bio *string
	if profile.Bio != "" {
		bio = &profile.Bio
	}

	return &models.UserProfile{
		User: models.User{
			ID:    profile.ID,
			Email: profile.Email,
			Name:  profile.Name,
		},
		// Extended fields from v1.6.0+
		Bio:      bio,
		City:     profile.City,
		Country:  profile.Country,
		Timezone: profile.Timezone,
		IsAdmin:  profile.IsAdmin,
	}
}

// FromEntityUserProfileList converts []entity.UserProfile to []*models.UserProfile
func FromEntityUserProfileList(profiles []*entity.UserProfile) []*models.UserProfile {
	if profiles == nil {
		return nil
	}

	result := make([]*models.UserProfile, 0, len(profiles))
	for _, profile := range profiles {
		result = append(result, FromEntityUserProfileToProfile(profile))
	}

	return result
}

// FromEntityRoleToModel converts entity.RoleInfo to models.Role
func FromEntityRoleToModel(role entity.RoleInfo) *models.Role {
	return &models.Role{
		Name:        role.Name,
		Description: role.Description,
	}
}

// FromEntityRoleList converts []entity.RoleInfo to []*models.Role
func FromEntityRoleList(roles []entity.RoleInfo) []*models.Role {
	if roles == nil {
		return nil
	}

	result := make([]*models.Role, 0, len(roles))
	for _, role := range roles {
		result = append(result, FromEntityRoleToModel(role))
	}

	return result
}
