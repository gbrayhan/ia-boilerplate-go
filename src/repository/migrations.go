package repository

import (
	"go.uber.org/zap"
	"os"
)

func (r *Repository) MigrateEntitiesGORM() error {
	if err := r.DB.AutoMigrate(&User{}, &RoleUser{}, &DeviceDetails{}, &Medicine{}, &ICDCie{}); err != nil {
		r.Logger.Error("Error migrating database entities", zap.Error(err))
		return err
	}
	r.Logger.Info("Database entities migrated successfully")

	if err := r.SeedInitialRole(); err != nil {
		r.Logger.Error("Error seeding initial role", zap.Error(err))
	}

	if err := r.SeedInitialUser(); err != nil {
		r.Logger.Error("Error seeding initial user", zap.Error(err))
		return err
	}
	r.Logger.Info("Seeding completed")
	return nil
}

func (r *Repository) SeedInitialRole() error {
	var count int64
	if err := r.DB.Model(&RoleUser{}).
		Where("name = ?", "admin").
		Count(&count).Error; err != nil {
		r.Logger.Error("Error checking existing role", zap.Error(err))
		return err
	}

	if count == 0 {
		role := RoleUser{Name: "admin", Description: "Administrator"}
		if err := r.DB.Create(&role).Error; err != nil {
			r.Logger.Error("Error creating initial role", zap.Error(err))
			return err
		}
		r.Logger.Info("Created initial role", zap.String("role", role.Name))
	} else {
		r.Logger.Debug("Admin role already exists", zap.Int64("count", count))
	}

	return nil
}

func (r *Repository) SeedInitialUser() error {
	email := os.Getenv("START_USER_EMAIL")
	pw := os.Getenv("START_USER_PW")
	if email == "" || pw == "" {
		r.Logger.Warn("Initial user seed skipped: START_USER_EMAIL or START_USER_PW not set")
		return nil
	}

	var count int64
	if err := r.DB.Model(&User{}).
		Where("email = ?", email).
		Count(&count).Error; err != nil {
		r.Logger.Error("Error checking existing user", zap.Error(err))
		return err
	}

	if count == 0 {
		hashed, err := r.Infrastructure.HashPassword(pw)
		if err != nil {
			r.Logger.Error("Error hashing initial user password", zap.Error(err))
			return err
		}
		var role RoleUser
		if err := r.DB.Where("name = ?", "admin").First(&role).Error; err != nil {
			r.Logger.Error("Error retrieving admin role for initial user", zap.Error(err))
			return err
		}

		user := User{
			Email:        email,
			HashPassword: hashed,
			FirstName:    "John",
			LastName:     "Doe",
			RoleID:       role.ID,
			Enabled:      true,
			JobPosition:  "Administrator",
		}
		if err := r.DB.Create(&user).Error; err != nil {
			r.Logger.Error("Error creating initial user", zap.Error(err))
			return err
		}
		r.Logger.Info("Created initial user", zap.String("email", email))
	} else {
		r.Logger.Debug("Initial user already exists", zap.String("email", email), zap.Int64("count", count))
	}

	return nil
}
