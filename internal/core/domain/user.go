package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	AdminRole  Role = "ADMIN"
	LesseeRole Role = "LESSEE"
	LessorRole Role = "LESSOR"
)

type Lifestyle string

const (
	Active   Lifestyle = "Active"
	Creative Lifestyle = "Creative"
	Social   Lifestyle = "Social"
	Relaxed  Lifestyle = "Relaxed"

	Football      Lifestyle = "Football"
	Basketball    Lifestyle = "Basketball"
	Tennis        Lifestyle = "Tennis"
	Swimming      Lifestyle = "Swimming"
	Running       Lifestyle = "Running"
	Cycling       Lifestyle = "Cycling"
	Badminton     Lifestyle = "Badminton"
	Yoga          Lifestyle = "Yoga"
	GymAndFitness Lifestyle = "Gym & Fitness"

	Music          Lifestyle = "Music"
	Dancing        Lifestyle = "Dancing"
	Photography    Lifestyle = "Photography"
	Painting       Lifestyle = "Painting"
	Gaming         Lifestyle = "Gaming"
	Reading        Lifestyle = "Reading"
	Writing        Lifestyle = "Writing"
	DIYAndCrafting Lifestyle = "DIY & Crafting"
	Cooking        Lifestyle = "Cooking"

	Extrovert Lifestyle = "Extrovert"
	Introvert Lifestyle = "Introvert"
	NightOwl  Lifestyle = "Night Owl"
	EarlyBird Lifestyle = "Early Bird"

	Traveler    Lifestyle = "Traveler"
	Backpacker  Lifestyle = "Backpacker"
	NatureLover Lifestyle = "Nature Lover"
	Camping     Lifestyle = "Camping"
	BeachLover  Lifestyle = "Beach Lover"

	DogLover Lifestyle = "Dog Lover"
	CatLover Lifestyle = "Cat Lover"

	Freelancer   Lifestyle = "Freelancer"
	Entrepreneur Lifestyle = "Entrepreneur"
	OfficeWorker Lifestyle = "Office Worker"
	RemoteWorker Lifestyle = "Remote Worker"
	Student      Lifestyle = "Student"
	SelfEmployed Lifestyle = "Self-Employed"
)

type User struct {
	ID                 uuid.UUID  `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreateAt           time.Time  `json:"createAt" gorm:"autoCreateTime"`
	UpdateAt           time.Time  `json:"updateAt" gorm:"autoUpdateTime"`
	Username           string     `json:"username" gorm:"unique" validate:"required"`
	Password           string     `json:"-" validate:"required,min=8"`
	Email              string     `json:"email" gorm:"unique" validate:"required,email"`
	Firstname          string     `json:"firstname"`
	Lastname           string     `json:"lastname"`
	NationalID         string     `json:"nationalID" `
	Gender             string     `json:"gender"`
	BirthDate          time.Time  `json:"birthDate" gorm:"type:DATE;default:null"`
	IsVerified         bool       `json:"isVerified" gorm:"default:false"`
	Role               *Role      `json:"role" gorm:"default:null"`
	FilledPersonalInfo bool       `json:"filledPersonalInfo" gorm:"default:false"`
	Lifestyle1         *Lifestyle `json:"lifestyle1" gorm:"default:null"`
	Lifestyle2         *Lifestyle `json:"lifestyle2" gorm:"default:null"`
	Lifestyle3         *Lifestyle `json:"lifestyle3" gorm:"default:null"`

	// studentEvidence
	StudentEvidence   string `json:"studentEvidence"`
	IsStudentVerified bool   `json:"isStudentVerified" gorm:"default:false"`
}
