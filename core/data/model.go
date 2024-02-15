package data

import "time"

type User struct {
	ID        uint      // Standard field for the primary key
	Name      string    // A regular string field
	Dailies   []Daily   `gorm:"many2many:user_dailies;"` // A many-to-many relationship
	CreatedAt time.Time // Automatically managed by GORM for creation time
	UpdatedAt time.Time // Automatically managed by GORM for update time
}

type Daily struct {
	ID        uint      // Standard field for the primary key
	Users     []User    `gorm:"many2many:user_dailies;"` // A many-to-many relationship
	CreatedAt time.Time // Automatically managed by GORM for creation time
	UpdatedAt time.Time // Automatically managed by GORM for update time
}

type Timer struct {
	ID      uint    // Standard field for the primary key
	UserID  uint    // The foreign key for the user
	DailyID uint    // The foreign key for the daily
	Elapsed float64 // The time elapsed
}
