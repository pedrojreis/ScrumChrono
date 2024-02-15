package data

import (
	"time"

	"gorm.io/driver/sqlite" // Sqlite driver based on CGO
	"gorm.io/gorm"
)

var db, err = gorm.Open(sqlite.Open("scrumchrono.db"), &gorm.Config{})

func InitDB(names []string) {
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&Daily{})
	db.AutoMigrate(&Timer{})

	//Create users if needed be
	for _, name := range names {
		var user User
		db.Where("name = ?", name).First(&user)
		if user.ID == 0 {
			db.Create(&User{Name: name})
		}
	}
}

func SaveMeeting(timers map[string]time.Duration) {
	var users []User
	var dailyTimers []Timer

	for name, duration := range timers {
		var user User
		db.Where("name = ?", name).First(&user)
		users = append(users, user)
		dailyTimers = append(dailyTimers, Timer{UserID: user.ID, Elapsed: duration.Seconds()})
	}

	daily := Daily{Users: users}
	db.Create(&daily)

	// I hate this second loop. There must be a better way.
	for i := range dailyTimers {
		dailyTimers[i].DailyID = daily.ID
	}

	db.Create(&dailyTimers)

	db.Commit()
}

func GetUserTimers(name string) []float64 {
	var user User
	db.Where("name = ?", name).First(&user)

	var timers []Timer
	db.Where("user_id = ?", user.ID).Find(&timers)

	// Convert timers to []float64
	var elapsedTimes []float64
	for _, timer := range timers {
		elapsedTimes = append(elapsedTimes, timer.Elapsed)
	}

	return elapsedTimes
}
