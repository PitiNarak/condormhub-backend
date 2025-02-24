package domain

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"
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

// Set of valid lifestyle values (for quick lookup)
var validLifestyles = map[Lifestyle]struct{}{
	Active:   {},
	Creative: {},
	Social:   {},
	Relaxed:  {},

	Football:      {},
	Basketball:    {},
	Tennis:        {},
	Swimming:      {},
	Running:       {},
	Cycling:       {},
	Badminton:     {},
	Yoga:          {},
	GymAndFitness: {},

	Music:          {},
	Dancing:        {},
	Photography:    {},
	Painting:       {},
	Gaming:         {},
	Reading:        {},
	Writing:        {},
	DIYAndCrafting: {},
	Cooking:        {},

	Extrovert: {},
	Introvert: {},
	NightOwl:  {},
	EarlyBird: {},

	Traveler:    {},
	Backpacker:  {},
	NatureLover: {},
	Camping:     {},
	BeachLover:  {},

	DogLover: {},
	CatLover: {},

	Freelancer:   {},
	Entrepreneur: {},
	OfficeWorker: {},
	RemoteWorker: {},
	Student:      {},
	SelfEmployed: {},
}

// IsValid checks if the lifestyle value is in the valid set.
func (l Lifestyle) IsValid() bool {
	_, exists := validLifestyles[l]
	return exists
}

type LifestyleArray []Lifestyle

// Scan converts PostgreSQL data to LifestyleArray
func (l *LifestyleArray) Scan(value interface{}) error {
	// Handle the empty array case from PostgreSQL ("{}")
	if value == "{}" {
		*l = []Lifestyle{}
		return nil
	}

	// Handle the case where the value is a string (PostgreSQL array format)
	str, ok := value.(string)
	if ok {
		// Remove surrounding curly braces
		str = str[1 : len(str)-1] // Removes `{}`

		// Use regex to correctly split elements handling quotes
		re := regexp.MustCompile(`"([^"]*)"|([^,]+)`)
		matches := re.FindAllStringSubmatch(str, -1)

		var lifestyles []Lifestyle
		for _, match := range matches {
			if match[1] != "" {
				lifestyles = append(lifestyles, Lifestyle(match[1])) // Remove surrounding quotes
			} else {
				lifestyles = append(lifestyles, Lifestyle(match[2]))
			}
		}

		*l = lifestyles
		return nil
	}

	// Handle the case where the value is a byte slice (PostgreSQL array format)
	byteArray, ok := value.([]byte)
	if ok {
		// Convert byte array to string and process similarly
		str := string(byteArray)
		str = str[1 : len(str)-1] // Removing the {} characters
		elements := strings.Split(str, ",")
		var lifestyles []Lifestyle
		for _, elem := range elements {
			elem = strings.TrimSpace(elem)
			lifestyles = append(lifestyles, Lifestyle(elem))
		}
		*l = lifestyles
		return nil
	}

	// If the value is neither a string nor a byte slice, return an error
	return fmt.Errorf("failed to scan LifestyleArray: %v", value)
}

// Value converts LifestyleArray to database format
func (l LifestyleArray) Value() (driver.Value, error) {
	if len(l) == 0 {
		return "{}", nil // Empty PostgreSQL array
	}

	// jsonData, err := json.Marshal(l)
	// fmt.Println(jsonData)

	strLifestyles := make([]string, len(l))
	for i, lifestyle := range l {
		strLifestyles[i] = string(lifestyle)
	}

	// Return the PostgreSQL array format without any extra string ("Lifestyles: ")
	return "{" + strings.Join(strLifestyles, ",") + "}", nil

	// if err != nil {
	// 	return nil, err
	// }
	// return string(jsonData), nil
}

type User struct {
	ID                 uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CreateAt           time.Time      `json:"createAt" gorm:"autoCreateTime"`
	UpdateAt           time.Time      `json:"updateAt" gorm:"autoUpdateTime"`
	Username           string         `json:"username" gorm:"unique" validate:"required"`
	Password           string         `json:"-" validate:"required,min=8"`
	Email              string         `json:"email" gorm:"unique" validate:"required,email"`
	Firstname          string         `json:"firstname"`
	Lastname           string         `json:"lastname"`
	NationalID         string         `json:"nationalID" `
	Gender             string         `json:"gender"`
	BirthDate          time.Time      `json:"birthDate" gorm:"type:DATE;default:null"`
	IsVerified         bool           `json:"isVerified" gorm:"default:false"`
	Role               *Role          `json:"role" gorm:"default:null"`
	FilledPersonalInfo bool           `json:"filledPersonalInfo" gorm:"default:false"`
	Lifestyles         LifestyleArray `json:"lifestyles" validate:"lifestyle" gorm:"type:lifestyle_tag[]"`

	// studentEvidence
	StudentEvidence   string `json:"studentEvidence"`
	IsStudentVerified bool   `json:"isStudentVerified" gorm:"default:false" `
}
