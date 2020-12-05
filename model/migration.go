package model

// migration auto migration
func migration() {
	if err := DB.AutoMigrate(&User{},
		&Card{}); err != nil {
		panic("auto migration err")
	}
}
