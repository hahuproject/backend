package auth_domain

import (
	"time"
)

type Course struct {
	ID       string `json:"courseId"`
	Name     string `json:"name"`
	CreditHr int    `json:"creditHr"`
	Color    string `json:"color"`
}

type Section struct {
	ID   string `json:"sectionId"`
	Name string `json:"name"`
	Year string `json:"year"`
}

type Stream struct {
	ID         string     `json:"classId"`
	Name       string     `json:"name"`
	Department Department `json:"department"`
}

type Department struct {
	ID   string `json:"departmentId"`
	Name string `json:"name"`
}

type User struct {
	ID         string    `json:"userId"`
	FirstName  string    `json:"firstName"`
	LastName   string    `json:"lastName"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	ProfilePic string    `json:"profilePic"`
	Verified   bool      `json:"verified"`
	Address    Address   `json:"address"`
	Type       string    `json:"type"`
	CreatedAt  time.Time `json:"createdAt"`
	Banned     bool      `json:"banned"`
}

type Admin struct {
	ID   string `json:"userId"`
	User User   `json:"user"`
}

type RegistryOfficer struct {
	ID   string `json:"userId"`
	User User   `json:"user"`
}

type SubRegistryOfficer struct {
	ID         string     `json:"userId"`
	User       User       `json:"user"`
	Department Department `json:"department"`
}

type DepartmentHead struct {
	ID   string `json:"userId"`
	User User   `json:"user"`
}

type Instructor struct {
	ID             string  `json:"userId"`
	SalaryRate     float32 `json:"salaryRate"`
	EmploymentType string  `json:"employmentType"`
	User           User    `json:"user"`
}

type Supervisor struct {
	ID       string `json:"userId"`
	Industry string `json:"industry"`
	User     User   `json:"user"`
}

type Student struct {
	ID                       string     `json:"userId"`
	MiddleName               string     `json:"middleName"`
	Gender                   string     `json:"gender"`
	Disablility              string     `json:"disability"`
	BirthDate                string     `json:"birthDate"`
	BirthPlace               string     `json:"birthPlace"`
	PreviousSchool           string     `json:"previousSchool"`
	Program                  string     `json:"program"`
	MatricResult             float32    `json:"matricResult"`
	AveragemarkForHighSchool float32    `json:"averageMarkForHighschool"`
	EmergencyContactName     string     `json:"emergencyContactName"`
	EmergencyContactPhone    string     `json:"emergencyContactPhone"`
	EmergencyContactRelation string     `json:"emergencyContactRelation"`
	EmergencyContactAddress  Address    `json:"address"`
	User                     User       `json:"user"`
	Department               Department `json:"department"`
	Stream                   Stream     `json:"stream"`
	Paid                     bool       `json:"paid"`
	Section                  Section    `json:"section"`
	Courses                  []Course   `json:"courses"`
}

// EntryType                string  `json:"entryType"`
