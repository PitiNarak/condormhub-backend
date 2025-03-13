package domain

import (
	"database/sql/driver"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/PitiNarak/condormhub-backend/internal/dto"
	"github.com/google/uuid"
)

type Role string

const (
	AdminRole  Role = "ADMIN"
	LesseeRole Role = "LESSEE"
	LessorRole Role = "LESSOR"
)

type Lifestyle string

var validLifestyles = []Lifestyle{
	"Active",
	"Creative",
	"Social",
	"Relaxed",

	"Football",
	"Basketball",
	"Tennis",
	"Swimming",
	"Running",
	"Cycling",
	"Badminton",
	"Yoga",
	"Gym & Fitness",

	"Music",
	"Dancing",
	"Photography",
	"Painting",
	"Gaming",
	"Reading",
	"Writing",
	"DIY & Crafting",
	"Cooking",

	"Extrovert",
	"Introvert",
	"Night Owl",
	"Early Bird",

	"Traveler",
	"Backpacker",
	"Nature Lover",
	"Camping",
	"Beach Lover",

	"Dog Lover",
	"Cat Lover",

	"Freelancer",
	"Entrepreneur",
	"Office Worker",
	"Remote Worker",
	"Student",
	"Self-Employed",
}

// IsValid checks if the lifestyle value is in the valid slice.
func (l Lifestyle) IsValid() bool {
	for _, validLifestyle := range validLifestyles {
		if l == validLifestyle {
			return true
		}
	}
	return false
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
	Role               Role           `json:"role" gorm:"default:null"`
	FilledPersonalInfo bool           `json:"filledPersonalInfo" gorm:"default:false"`
	Lifestyles         LifestyleArray `json:"lifestyles" validate:"lifestyle" gorm:"type:lifestyle_tag[]"`
	PhoneNumber        string         `json:"phoneNumber"`
	// studentEvidence
	StudentEvidence   string `json:"studentEvidence"`
	IsStudentVerified bool   `json:"isStudentVerified" gorm:"default:false" `
}

func (u *User) ToDTO() dto.UserResponse {
	lifestyles := make([]string, len(u.Lifestyles))
	for i, v := range u.Lifestyles {
		lifestyles[i] = string(v)
	}
	return dto.UserResponse{
		ID:                 u.ID,
		Username:           u.Username,
		Email:              u.Email,
		Firstname:          u.Firstname,
		Lastname:           u.Lastname,
		Gender:             u.Gender,
		BirthDate:          u.BirthDate,
		IsVerified:         u.IsVerified,
		Role:               string(u.Role),
		FilledPersonalInfo: u.FilledPersonalInfo,
		Lifestyles:         lifestyles,
		PhoneNumber:        u.PhoneNumber,
		StudentEvidence:    u.StudentEvidence,
		IsStudentVerified:  u.IsStudentVerified,
	}
}
