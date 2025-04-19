package db

import (
	"fmt"
	"ia-boilerplate/utils"
	"os"
)

func MigrateEntitiesGORM() error {
	err := DB.AutoMigrate(&User{}, &RoleUser{}, &DeviceDetails{}, &Medicine{}, &ICDCie{})
	if err != nil {
		fmt.Println("Error migrating the database:", err)
		return err
	}

	err = SeedInitialUser()
	if err != nil {
		fmt.Println("Error seeding initial user:", err)
		return err
	}
	err = SeedInitialRole()
	if err != nil {
		fmt.Println("Error seeding initial role:", err)
	}
	return err
}
func SeedInitialRole() error {
	var count int64
	if err := DB.Model(&RoleUser{}).
		Where("name = ?", "admin").
		Count(&count).Error; err != nil {
		return fmt.Errorf("error checking existing role: %w", err)
	}

	if count == 0 {
		role := RoleUser{
			Name:        "admin",
			Description: "Administrator",
		}
		if err := DB.Create(&role).Error; err != nil {
			return fmt.Errorf("error creating initial role: %w", err)
		}
		fmt.Printf("✔️  Created initial role %s\n", role.Name)
	}

	return nil
}

func SeedInitialUser() error {
	email := os.Getenv("START_USER_EMAIL")
	pw := os.Getenv("START_USER_PW")
	if email == "" || pw == "" {
		return nil
	}

	var count int64
	if err := DB.Model(&User{}).
		Where("email = ?", email).
		Count(&count).Error; err != nil {
		return fmt.Errorf("error checking existing user: %w", err)
	}

	if count == 0 {
		hashed, err := utils.HashPassword(pw)
		if err != nil {
			return fmt.Errorf("error hashing password: %w", err)
		}

		user := User{
			Email:        email,
			HashPassword: hashed,
			FirstName:    "John",
			LastName:     "Doe",
			Enabled:      true,
			JobPosition:  "Administrator",
		}
		if err := DB.Create(&user).Error; err != nil {
			return fmt.Errorf("error creating initial user: %w", err)
		}
		fmt.Printf("✔️  Created initial user %s\n", email)
	}

	return nil
}
