package auth_service

import (
	"log"
	"time"

	auth_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/domain"
	authRepo "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/repo"
	authError "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/utils"
	"github.com/dgrijalva/jwt-go/v4"
	"golang.org/x/crypto/bcrypt"
)

type AuthServicePort interface {
	RegisterAdmin(user auth_domain.Admin) (auth_domain.Admin, error)
	RegisterRegistryOfficer(user auth_domain.RegistryOfficer) (auth_domain.RegistryOfficer, error)
	RegisterSubRegistryOfficer(user auth_domain.SubRegistryOfficer) (auth_domain.SubRegistryOfficer, error)
	RegisterDepartmentHead(user auth_domain.DepartmentHead) (auth_domain.DepartmentHead, error)
	RegisterInstructor(user auth_domain.Instructor) (auth_domain.Instructor, error)
	RegisterSupervisor(user auth_domain.Supervisor) (auth_domain.Supervisor, error)
	RegisterStudent(token string, user auth_domain.Student) (auth_domain.Student, error)
	Login(username, password string) (auth_domain.User, string, error)
	GetMe(token string) (auth_domain.User, error)
	GetRole(id string) (string, error)
	ApproveUser(token string, id string) error
	BanUser(token, id string) error
	UnBanUser(token, id string) error
	GetUsers(token string) ([]auth_domain.User, error)
	GetStudent(token string) (auth_domain.Student, error)
	GetInstructor(token string) (auth_domain.Instructor, error)
	GetSupervisor(token string) (auth_domain.Supervisor, error)
	ChangePassword(token, oldPassword, newPassword string) error
	UpdateProfile(token string, user auth_domain.User) error

	//Students
	GetStudents(token string) ([]auth_domain.Student, error)
	UpdateStudentPayment(token, userId string, status bool) (auth_domain.Student, error)
}

type AuthServiceAdapter struct {
	log  *log.Logger
	repo authRepo.AuthRepositoryPort
}

func NewAuthService(log *log.Logger, repo authRepo.AuthRepositoryPort) AuthServicePort {
	return &AuthServiceAdapter{
		log, repo,
	}
}

func (authService AuthServiceAdapter) RegisterAdmin(user auth_domain.Admin) (auth_domain.Admin, error) {

	var addedUser auth_domain.Admin

	//Check Data Validity
	if user.User.FirstName == "" || len(user.User.FirstName) < 2 {
		return addedUser, authError.ErrInvalidFirstName
	}

	if user.User.LastName == "" || len(user.User.LastName) < 2 {
		return addedUser, authError.ErrInvalidLastName
	}

	if user.User.Email == "" {
		return addedUser, authError.ErrInvalidEmail
	}

	if user.User.Phone == "" || len(user.User.Phone) != 10 {
		return addedUser, authError.ErrInvalidPhone
	}

	if user.User.Username == "" || len(user.User.Username) < 5 {
		return addedUser, authError.ErrInvalidUsername
	}

	if user.User.Password == "" || len(user.User.Password) < 6 {
		return addedUser, authError.ErrInvalidPassword
	}

	addedUser, err := authService.repo.StoreAdmin(user)
	if err != nil {
		return addedUser, err
	}
	return addedUser, nil
}

func (authService AuthServiceAdapter) RegisterSubRegistryOfficer(user auth_domain.SubRegistryOfficer) (auth_domain.SubRegistryOfficer, error) {

	var addedUser auth_domain.SubRegistryOfficer

	//Check Data Validity
	if user.User.FirstName == "" || len(user.User.FirstName) < 2 {
		return addedUser, authError.ErrInvalidFirstName
	}

	if user.User.LastName == "" || len(user.User.LastName) < 2 {
		return addedUser, authError.ErrInvalidLastName
	}

	if user.User.Email == "" {
		return addedUser, authError.ErrInvalidEmail
	}

	if user.User.Phone == "" || len(user.User.Phone) != 10 {
		return addedUser, authError.ErrInvalidPhone
	}

	if user.User.Username == "" || len(user.User.Username) < 5 {
		return addedUser, authError.ErrInvalidUsername
	}

	if user.User.Password == "" || len(user.User.Password) < 6 {
		return addedUser, authError.ErrInvalidPassword
	}

	addedUser, err := authService.repo.StoreSubRegistryOfficer(user)
	if err != nil {
		return addedUser, err
	}
	return addedUser, nil
}

func (authService AuthServiceAdapter) RegisterRegistryOfficer(user auth_domain.RegistryOfficer) (auth_domain.RegistryOfficer, error) {

	var addedUser auth_domain.RegistryOfficer

	//Check Data Validity
	if user.User.FirstName == "" || len(user.User.FirstName) < 2 {
		return addedUser, authError.ErrInvalidFirstName
	}

	if user.User.LastName == "" || len(user.User.LastName) < 2 {
		return addedUser, authError.ErrInvalidLastName
	}

	if user.User.Email == "" {
		return addedUser, authError.ErrInvalidEmail
	}

	if user.User.Phone == "" || len(user.User.Phone) != 10 {
		return addedUser, authError.ErrInvalidPhone
	}

	if user.User.Username == "" || len(user.User.Username) < 5 {
		return addedUser, authError.ErrInvalidUsername
	}

	if user.User.Password == "" || len(user.User.Password) < 6 {
		return addedUser, authError.ErrInvalidPassword
	}

	addedUser, err := authService.repo.StoreRegistryOfficer(user)
	if err != nil {
		return addedUser, err
	}
	return addedUser, nil
}

func (authService AuthServiceAdapter) RegisterDepartmentHead(user auth_domain.DepartmentHead) (auth_domain.DepartmentHead, error) {
	var addedUser auth_domain.DepartmentHead

	//Check Data Validity
	if user.User.FirstName == "" || len(user.User.FirstName) < 2 {
		return addedUser, authError.ErrInvalidFirstName
	}

	if user.User.LastName == "" || len(user.User.LastName) < 2 {
		return addedUser, authError.ErrInvalidLastName
	}

	if user.User.Email == "" {
		return addedUser, authError.ErrInvalidEmail
	}

	if user.User.Phone == "" || len(user.User.Phone) != 10 {
		return addedUser, authError.ErrInvalidPhone
	}

	if user.User.Username == "" || len(user.User.Username) < 5 {
		return addedUser, authError.ErrInvalidUsername
	}

	if user.User.Password == "" || len(user.User.Password) < 6 {
		return addedUser, authError.ErrInvalidPassword
	}

	addedUser, err := authService.repo.StoreDepartmentHead(user)
	if err != nil {
		return addedUser, err
	}
	return addedUser, nil
}

func (authService AuthServiceAdapter) RegisterInstructor(user auth_domain.Instructor) (auth_domain.Instructor, error) {
	var addedUser auth_domain.Instructor

	//Check Data Validity
	if user.User.FirstName == "" || len(user.User.FirstName) < 2 {
		return addedUser, authError.ErrInvalidFirstName
	}

	if user.User.LastName == "" || len(user.User.LastName) < 2 {
		return addedUser, authError.ErrInvalidLastName
	}

	if user.User.Email == "" {
		return addedUser, authError.ErrInvalidEmail
	}

	if user.User.Phone == "" || len(user.User.Phone) != 10 {
		return addedUser, authError.ErrInvalidPhone
	}

	if user.User.Username == "" || len(user.User.Username) < 5 {
		return addedUser, authError.ErrInvalidUsername
	}

	if user.User.Password == "" || len(user.User.Password) < 6 {
		return addedUser, authError.ErrInvalidPassword
	}

	addedUser, err := authService.repo.StoreInstructor(user)
	if err != nil {
		return addedUser, err
	}
	return addedUser, nil
}

func (authService AuthServiceAdapter) RegisterSupervisor(user auth_domain.Supervisor) (auth_domain.Supervisor, error) {
	var addedUser auth_domain.Supervisor

	//Check Data Validity
	if user.User.FirstName == "" || len(user.User.FirstName) < 2 {
		return addedUser, authError.ErrInvalidFirstName
	}

	if user.User.LastName == "" || len(user.User.LastName) < 2 {
		return addedUser, authError.ErrInvalidLastName
	}

	if user.User.Email == "" {
		return addedUser, authError.ErrInvalidEmail
	}

	if user.User.Phone == "" || len(user.User.Phone) != 10 {
		return addedUser, authError.ErrInvalidPhone
	}

	if user.User.Username == "" || len(user.User.Username) < 5 {
		return addedUser, authError.ErrInvalidUsername
	}

	if user.User.Password == "" || len(user.User.Password) < 6 {
		return addedUser, authError.ErrInvalidPassword
	}

	addedUser, err := authService.repo.StoreSupervisor(user)
	if err != nil {
		return addedUser, err
	}
	return addedUser, nil
}

func (authService AuthServiceAdapter) RegisterStudent(token string, user auth_domain.Student) (auth_domain.Student, error) {

	var _user auth_domain.User
	if token != "" {
		_user, _ = authService.GetMe(token)
	}

	var addedUser auth_domain.Student

	//Check Data Validity
	if user.User.FirstName == "" || len(user.User.FirstName) < 2 {
		return addedUser, authError.ErrInvalidFirstName
	}

	if user.User.LastName == "" || len(user.User.LastName) < 2 {
		return addedUser, authError.ErrInvalidLastName
	}

	if user.User.Email == "" {
		return addedUser, authError.ErrInvalidEmail
	}

	if user.User.Phone == "" || len(user.User.Phone) != 10 {
		return addedUser, authError.ErrInvalidPhone
	}

	if user.User.Username == "" || len(user.User.Username) < 5 {
		return addedUser, authError.ErrInvalidUsername
	}

	if user.User.Password == "" || len(user.User.Password) < 6 {
		return addedUser, authError.ErrInvalidPassword
	}

	var verified bool = false
	if _user.Type == "REGISTRY_OFFICER" {
		verified = true
	}

	addedUser, err := authService.repo.StoreStudent(verified, user)
	if err != nil {
		return addedUser, err
	}
	return addedUser, nil
}

func (authService AuthServiceAdapter) Login(username, password string) (auth_domain.User, string, error) {
	user, err := authService.repo.FindUserByUsername(username)
	if err != nil {
		return auth_domain.User{}, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return auth_domain.User{}, "", authError.ErrIncorrectPassword
	}

	if !user.Verified {
		return auth_domain.User{}, "", authError.ErrUserNotVerified
	}

	if user.Banned {
		return auth_domain.User{}, "", authError.ErrUserBanned
	}

	// authService.log.Println("oooo 1")

	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["id"] = user.ID
	atClaims["exp"] = time.Now().Add(time.Hour * 24 * 360).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte("hahu_sms_secret"))
	if err != nil {
		// authService.log.Println("oooo 2")
		authService.log.Println(err)
		return auth_domain.User{}, "", err
	}

	// authService.log.Println("oooo 3")
	return user, token, nil
}

func (authService AuthServiceAdapter) GetMe(token string) (auth_domain.User, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("hahu_sms_secret"), nil
	})

	if err != nil {
		return auth_domain.User{}, authError.ErrUnauthorizedRequest
	}

	//authService.log.Println(claims["id"])
	if claims["id"] == nil || claims["id"] == "" {
		//authService.log.Println("ApproveUser SERVICE 2- " + claims["id"].(string))
		return auth_domain.User{}, authError.ErrUnauthorizedRequest
	}

	user, err := authService.repo.FindUserById(claims["id"].(string))
	if err != nil {
		return auth_domain.User{}, authError.ErrUserNotFound
	}

	return user, nil
}

func (authService AuthServiceAdapter) GetRole(id string) (string, error) {
	role, err := authService.repo.GetUserRoleByID(id)
	if err != nil {
		return "", err
	}
	return role, nil
}
func (authService AuthServiceAdapter) ApproveUser(token string, id string) error {

	user, err := authService.GetMe(token)
	if err != nil {
		return err
	}

	role, err := authService.GetRole(user.ID)
	if err != nil {
		//authService.log.Println("ApproveUser SERVICE 3- " + err.Error())
		return err
	}

	if role != "SUPER_ADMIN" && role != "ADMIN" && role != "REGISTRY_OFFICER" {
		//authService.log.Println("ApproveUser SERVICE 4- " + role)
		return authError.ErrUnauthorizedRequest
	}

	//authService.log.Println("ApproveUser SERVICE 5- ")
	err = authService.repo.ApproveUser(id)
	if err != nil {
		//authService.log.Println("ApproveUser SERVICE 5- " + err.Error())
		return err
	}

	// authService.repo.ApproveUser(id)
	return nil
}
func (authService AuthServiceAdapter) BanUser(token, id string) error {
	user, err := authService.GetMe(token)
	if err != nil {
		return err
	}

	role, err := authService.GetRole(user.ID)
	if err != nil {
		return err
	}

	if role != "SUPER_ADMIN" {
		return authError.ErrUnauthorizedRequest
	}

	err = authService.repo.BanUser(id)
	if err != nil {
		return err
	}

	return nil
}

func (authService AuthServiceAdapter) UnBanUser(token, id string) error {
	user, err := authService.GetMe(token)
	if err != nil {
		return err
	}

	role, err := authService.GetRole(user.ID)
	if err != nil {
		return err
	}

	if role != "SUPER_ADMIN" {
		return authError.ErrUnauthorizedRequest
	}

	err = authService.repo.UnBanUser(id)
	if err != nil {
		return err
	}

	return nil
}
func (authService AuthServiceAdapter) GetUsers(token string) ([]auth_domain.User, error) {

	var user auth_domain.User

	claims := jwt.MapClaims{}
	jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("hahu_sms_secret"), nil
	})

	if claims["id"] != nil || claims["id"] != "" {
		user, _ = authService.repo.FindUserById(claims["id"].(string))
	}

	users, err := authService.repo.FindAllUsers()
	if err != nil {
		return users, err
	}

	var filteredUsers []auth_domain.User = make([]auth_domain.User, 0)
	for i := 0; i < len(users); i++ {
		if users[i].ID != user.ID && users[i].Type != "SUPER_ADMIN" {
			filteredUsers = append(filteredUsers, users[i])
		}
	}

	return filteredUsers, nil
}

func (authService AuthServiceAdapter) GetStudent(token string) (auth_domain.Student, error) {
	var _student auth_domain.Student

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("hahu_sms_secret"), nil
	})

	if err != nil {
		return _student, authError.ErrUnauthorizedRequest
	}

	//authService.log.Println(claims["id"])
	if claims["id"] == nil || claims["id"] == "" {
		//authService.log.Println("ApproveUser SERVICE 2- " + claims["id"].(string))
		return _student, authError.ErrUnauthorizedRequest
	}

	user, err := authService.repo.FindUserById(claims["id"].(string))
	if err != nil {
		return _student, authError.ErrUserNotFound
	}

	return authService.repo.FindStudent(user.ID)
}
func (authService AuthServiceAdapter) GetInstructor(token string) (auth_domain.Instructor, error) {
	var _instructor auth_domain.Instructor

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("hahu_sms_secret"), nil
	})

	if err != nil {
		return _instructor, authError.ErrUnauthorizedRequest
	}

	//authService.log.Println(claims["id"])
	if claims["id"] == nil || claims["id"] == "" {
		//authService.log.Println("ApproveUser SERVICE 2- " + claims["id"].(string))
		return _instructor, authError.ErrUnauthorizedRequest
	}

	user, err := authService.repo.FindUserById(claims["id"].(string))
	if err != nil {
		return _instructor, authError.ErrUserNotFound
	}

	return authService.repo.FindInstructor(user.ID)
}

func (authService AuthServiceAdapter) GetSupervisor(token string) (auth_domain.Supervisor, error) {
	var _supervisor auth_domain.Supervisor

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("hahu_sms_secret"), nil
	})

	if err != nil {
		return _supervisor, authError.ErrUnauthorizedRequest
	}

	//authService.log.Println(claims["id"])
	if claims["id"] == nil || claims["id"] == "" {
		//authService.log.Println("ApproveUser SERVICE 2- " + claims["id"].(string))
		return _supervisor, authError.ErrUnauthorizedRequest
	}

	user, err := authService.repo.FindUserById(claims["id"].(string))
	if err != nil {
		return _supervisor, authError.ErrUserNotFound
	}

	return authService.repo.FindSupervisor(user.ID)
}

func (authService AuthServiceAdapter) ChangePassword(token, oldPassword, newPassword string) error {

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("hahu_sms_secret"), nil
	})

	if err != nil {
		return authError.ErrUnauthorizedRequest
	}

	if claims["id"] == nil || claims["id"] == "" {
		return authError.ErrUnauthorizedRequest
	}

	user, err := authService.repo.FindUserById(claims["id"].(string))
	if err != nil {
		return authError.ErrUserNotFound
	}

	print(oldPassword)
	print(user.Password)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return authError.ErrIncorrectPassword
	}

	return authService.repo.ChangePassword(user.ID, newPassword)
}

func (authService AuthServiceAdapter) UpdateProfile(token string, user auth_domain.User) error {

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("hahu_sms_secret"), nil
	})

	if err != nil {
		return authError.ErrUnauthorizedRequest
	}

	if claims["id"] == nil || claims["id"] == "" {
		return authError.ErrUnauthorizedRequest
	}

	_user, err := authService.repo.FindUserById(claims["id"].(string))
	if err != nil {
		return authError.ErrUserNotFound
	}

	user.ID = _user.ID

	return authService.repo.UpdateProfile(user)
}

func (service AuthServiceAdapter) GetStudents(token string) ([]auth_domain.Student, error) {

	var students []auth_domain.Student = make([]auth_domain.Student, 0)

	user, err := service.GetMe(token)
	if err != nil {
		return students, err
	}

	role, err := service.GetRole(user.ID)
	if err != nil {
		return students, err
	}

	if role != "SUPER_ADMIN" && role != "ADMIN" && role != "REGISTRY_OFFICER" && role != "DEPARTMENT_HEAD" {
		return students, authError.ErrUnauthorizedRequest
	}

	students, err = service.repo.FindStudents()
	if err != nil {
		return students, err
	}

	return students, nil
}
func (service AuthServiceAdapter) UpdateStudentPayment(token, userId string, status bool) (auth_domain.Student, error) {
	var updatedStudent auth_domain.Student

	user, err := service.GetMe(token)
	if err != nil {
		return updatedStudent, err
	}

	role, err := service.GetRole(user.ID)
	if err != nil {
		return updatedStudent, err
	}

	if role != "REGISTRY_OFFICER" {
		return updatedStudent, authError.ErrUnauthorizedRequest
	}

	updatedStudent, err = service.repo.UpdateStudentPaymentSatus(userId, status)
	if err != nil {
		return updatedStudent, err
	}

	return updatedStudent, nil
}
