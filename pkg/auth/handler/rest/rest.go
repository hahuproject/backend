package auth_rest_handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	auth_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/domain"
	authHandler "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/handler"
	authService "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/service"
	auth_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/utils"
)

type RestAuthHandlerAdapter struct {
	log         *log.Logger
	authService authService.AuthServicePort
}

func NewRestAuthHandlerAdapter(log *log.Logger, authService authService.AuthServicePort) authHandler.RestAuthHandlerPort {
	return &RestAuthHandlerAdapter{log: log, authService: authService}
}

func checkPOST(r *http.Request, w http.ResponseWriter) bool {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(r.Method + " Method not supported "))
		return false
	}
	return true
}

func populateUserFromFormValue(r *http.Request) auth_domain.User {
	var user auth_domain.User

	user.FirstName = r.FormValue("firstName")
	user.LastName = r.FormValue("lastName")
	user.Email = r.FormValue("email")
	user.Phone = r.FormValue("phone")
	user.Username = r.FormValue("username")
	user.Password = r.FormValue("password")
	user.Address.Country = r.FormValue("country")
	user.Address.Region = r.FormValue("region")
	user.Address.City = r.FormValue("city")
	user.Address.SubCity = r.FormValue("subCity")
	woreda, _ := strconv.Atoi(r.FormValue("woreda"))
	user.Address.Woreda = woreda
	user.Address.HouseNo = r.FormValue("houseNo")

	return user
}

func sendMasrshaledResponse(data interface{}, w http.ResponseWriter) {
	marshaledAddedUser, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshaledAddedUser)
}

func storeProfilePicFromFormFile(r *http.Request, user auth_domain.User) string {

	var storedFileAdd string = ""

	file, fileHandler, err := r.FormFile("profilePic")
	if err == nil {
		if file != nil {
			storedFile, err := auth_utils.StoreFileFromMultiPartRequest(file, user.Username+time.Now().UTC().String()+"."+strings.Split(fileHandler.Filename, ".")[len(strings.Split(fileHandler.Filename, "."))-1], "uploads/images/profile")
			if err != nil {
				return storedFileAdd
			}

			storedFileAdd = r.Host + strings.Replace(storedFile.Name(), ".", "", 1)
		}

		defer file.Close()
	}

	return storedFileAdd
}

func (restAuthHandler *RestAuthHandlerAdapter) GetRegisterAdmin(w http.ResponseWriter, r *http.Request) {
	check := checkPOST(r, w)
	if !check {
		return
	}

	r.ParseMultipartForm(32 << 20)

	var user auth_domain.Admin

	user.User = populateUserFromFormValue(r)

	//Store file using - uploads/images/profile/username + time + extension

	user.User.ProfilePic = auth_utils.StoreProfilePicFromFormFileToAWSS3(r, user.User)

	addedUser, err := restAuthHandler.authService.RegisterAdmin(user)
	if err != nil {
		w.WriteHeader(http.StatusFailedDependency)
		w.Write([]byte(err.Error()))
		return
	}

	sendMasrshaledResponse(addedUser, w)

}

func (restAuthHandler *RestAuthHandlerAdapter) GetRegisterRegistryOfficer(w http.ResponseWriter, r *http.Request) {
	check := checkPOST(r, w)
	if !check {
		return
	}

	r.ParseMultipartForm(32 << 20)

	var user auth_domain.RegistryOfficer

	user.User = populateUserFromFormValue(r)

	//Store file using - uploads/images/profile/username + time + extension

	user.User.ProfilePic = auth_utils.StoreProfilePicFromFormFileToAWSS3(r, user.User)

	addedUser, err := restAuthHandler.authService.RegisterRegistryOfficer(user)
	if err != nil {
		w.WriteHeader(http.StatusFailedDependency)
		w.Write([]byte(err.Error()))
		return
	}

	sendMasrshaledResponse(addedUser, w)
}

func (restAuthHandler *RestAuthHandlerAdapter) GetRegisterSubRegistryOfficer(w http.ResponseWriter, r *http.Request) {
	check := checkPOST(r, w)
	if !check {
		return
	}

	r.ParseMultipartForm(32 << 20)

	var user auth_domain.SubRegistryOfficer

	user.User = populateUserFromFormValue(r)
	user.Department.ID = r.FormValue("departmentId")

	//Store file using - uploads/images/profile/username + time + extension

	user.User.ProfilePic = auth_utils.StoreProfilePicFromFormFileToAWSS3(r, user.User)

	addedUser, err := restAuthHandler.authService.RegisterSubRegistryOfficer(user)
	if err != nil {
		w.WriteHeader(http.StatusFailedDependency)
		w.Write([]byte(err.Error()))
		return
	}

	sendMasrshaledResponse(addedUser, w)
}

func (restAuthHandler *RestAuthHandlerAdapter) GetRegisterDepartmentHead(w http.ResponseWriter, r *http.Request) {
	check := checkPOST(r, w)
	if !check {
		return
	}

	r.ParseMultipartForm(32 << 20)

	var user auth_domain.DepartmentHead

	user.User = populateUserFromFormValue(r)

	//Store file using - uploads/images/profile/username + time + extension

	user.User.ProfilePic = auth_utils.StoreProfilePicFromFormFileToAWSS3(r, user.User)

	addedUser, err := restAuthHandler.authService.RegisterDepartmentHead(user)
	if err != nil {
		w.WriteHeader(http.StatusFailedDependency)
		w.Write([]byte(err.Error()))
		return
	}

	sendMasrshaledResponse(addedUser, w)
}
func (restAuthHandler *RestAuthHandlerAdapter) GetRegisterInstructor(w http.ResponseWriter, r *http.Request) {
	check := checkPOST(r, w)
	if !check {
		return
	}
	print("check 0")

	r.ParseMultipartForm(32 << 20)
	print("check 1")

	var user auth_domain.Instructor
	print("check 2")

	user.User = populateUserFromFormValue(r)
	print("check 3")
	user.EmploymentType = r.FormValue("employmentType")
	print("check 4")

	//Store file using - uploads/images/profile/username + time + extension

	user.User.ProfilePic = auth_utils.StoreProfilePicFromFormFileToAWSS3(r, user.User)
	print("check 5")

	addedUser, err := restAuthHandler.authService.RegisterInstructor(user)
	print("check 6")
	if err != nil {
		print("check 7")
		w.WriteHeader(http.StatusFailedDependency)
		w.Write([]byte(err.Error()))
		return
	}

	print("check 8")
	sendMasrshaledResponse(addedUser, w)
	print("check 9")
}
func (restAuthHandler *RestAuthHandlerAdapter) GetRegisterSupervisor(w http.ResponseWriter, r *http.Request) {
	check := checkPOST(r, w)
	if !check {
		return
	}

	r.ParseMultipartForm(32 << 20)

	var user auth_domain.Supervisor

	user.User = populateUserFromFormValue(r)
	user.Industry = r.FormValue("industry")

	//Store file using - uploads/images/profile/username + time + extension

	user.User.ProfilePic = auth_utils.StoreProfilePicFromFormFileToAWSS3(r, user.User)

	addedUser, err := restAuthHandler.authService.RegisterSupervisor(user)
	if err != nil {
		w.WriteHeader(http.StatusFailedDependency)
		w.Write([]byte(err.Error()))
		return
	}

	sendMasrshaledResponse(addedUser, w)
}
func (restAuthHandler *RestAuthHandlerAdapter) GetRegisterStudent(w http.ResponseWriter, r *http.Request) {
	check := checkPOST(r, w)
	if !check {
		return
	}

	r.ParseMultipartForm(32 << 20)

	var user auth_domain.Student

	user.User = populateUserFromFormValue(r)

	user.MiddleName = r.FormValue("middleName")
	user.Gender = r.FormValue("gender")
	user.BirthDate = r.FormValue("birthDate")
	user.BirthPlace = r.FormValue("birthPlace")
	user.Disablility = r.FormValue("disability")

	user.PreviousSchool = r.FormValue("previousSchool")
	averageMarkForHighschool, _ := strconv.Atoi(r.FormValue("avarageMarkForHighSchool"))
	user.AveragemarkForHighSchool = float32(averageMarkForHighschool)
	matricResult, _ := strconv.Atoi(r.FormValue("matricResult"))
	user.MatricResult = float32(matricResult)
	user.Program = r.FormValue("program")
	user.Stream.ID = r.FormValue("streamId")
	user.Department.ID = r.FormValue("departmentId")
	user.Department.Name = r.FormValue("departmentName")
	user.EmergencyContactName = r.FormValue("emergenctContactName")
	user.EmergencyContactPhone = r.FormValue("emergenctContactPhone")
	user.EmergencyContactRelation = r.FormValue("emergenctContactRelation")
	user.EmergencyContactAddress.Country = r.FormValue("emergenctContactCountry")
	user.EmergencyContactAddress.Region = r.FormValue("emergenctContactRegion")
	user.EmergencyContactAddress.City = r.FormValue("emergenctContactCity")
	user.EmergencyContactAddress.SubCity = r.FormValue("emergenctContactSubcity")
	emergenctContactWoreda, err := strconv.Atoi(r.FormValue("emergenctContactWoreda"))
	if err == nil {
		user.EmergencyContactAddress.Woreda = emergenctContactWoreda

	}
	user.EmergencyContactAddress.HouseNo = r.FormValue("emergenctContactHouseno")

	//Store file using - uploads/images/profile/username + time + extension

	user.User.ProfilePic = auth_utils.StoreProfilePicFromFormFileToAWSS3(r, user.User)

	var _token string

	if r.Header.Get("Authorization") != "" && len(strings.Split(r.Header.Get("Authorization"), " ")) > 0 {
		_token = strings.Split(r.Header.Get("Authorization"), " ")[1]
	}

	addedUser, err := restAuthHandler.authService.RegisterStudent(_token, user)
	if err != nil {
		w.WriteHeader(http.StatusFailedDependency)
		w.Write([]byte(err.Error()))
		return
	}

	sendMasrshaledResponse(addedUser, w)
}
func (restAuthHandler *RestAuthHandlerAdapter) GetLogin(w http.ResponseWriter, r *http.Request) {

	type LoginCredential struct {
		Username string
		Password string
	}

	var loginCredentials LoginCredential
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&loginCredentials)
	user, token, err := restAuthHandler.authService.Login(loginCredentials.Username, loginCredentials.Password)
	if err != nil {
		w.WriteHeader(http.StatusFailedDependency)
		w.Write([]byte(err.Error()))
		return
	}

	type LoginReturn struct {
		Token string `json:"token"`
		ID    string `json:"id"`
	}

	res_ := LoginReturn{Token: token, ID: user.ID}

	res, _ := json.Marshal(res_)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}

func (restAuthHandler RestAuthHandlerAdapter) GetUnBanUser(w http.ResponseWriter, r *http.Request) {
	err := restAuthHandler.authService.UnBanUser(strings.Split(r.Header.Get("Authorization"), " ")[1], r.URL.Query()["id"][0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	// w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully unbanned user"))
}

func (restAuthHandler RestAuthHandlerAdapter) GetBanUser(w http.ResponseWriter, r *http.Request) {
	err := restAuthHandler.authService.BanUser(strings.Split(r.Header.Get("Authorization"), " ")[1], r.URL.Query()["id"][0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	// w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully banned user"))
}

func (restAuthHandler RestAuthHandlerAdapter) GetApproveUser(w http.ResponseWriter, r *http.Request) {
	err := restAuthHandler.authService.ApproveUser(strings.Split(r.Header.Get("Authorization"), " ")[1], r.URL.Query()["id"][0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	// w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully approved user"))
}

func (restAuthHandler RestAuthHandlerAdapter) GetMe(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.Header.Get("Authorization"), " ")
	if token == nil || len(token) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not Authorized"))
		return
	}

	_userType := r.URL.Query().Get("type")

	var user interface{}
	var err error

	switch _userType {
	case "":
		{
			user, err = restAuthHandler.authService.GetMe(token[1])
			break
		}
	case "student":
		{
			user, err = restAuthHandler.authService.GetStudent(token[1])
			break
		}
	case "instructor":
		{
			user, err = restAuthHandler.authService.GetInstructor(token[1])
			break
		}
	case "supervisor":
		{
			user, err = restAuthHandler.authService.GetSupervisor(token[1])
			break
		}

	}

	// restAuthHandler.log.Println("GetMe called", token[1])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	// w.WriteHeader(http.StatusOK)
	marshaledUser, err := json.Marshal(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshaledUser)
}

func (restAuthHandler RestAuthHandlerAdapter) GetUsers(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.Header.Get("Authorization"), " ")
	if token == nil || len(token) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not Authorized"))
		return
	}

	users, err := restAuthHandler.authService.GetUsers(token[1])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	marshaledUsers, err := json.Marshal(users)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(marshaledUsers)
}

func (restAuthHandler RestAuthHandlerAdapter) GetChangePassword(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.Header.Get("Authorization"), " ")
	if token == nil || len(token) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not Authorized"))
		return
	}

	type PassUpdateType struct {
		OldPassword string `json:"oldPassword"`
		NewPassword string `json:"newPassword"`
	}

	var _passUpdate PassUpdateType

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&_passUpdate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	defer r.Body.Close()

	err = restAuthHandler.authService.ChangePassword(token[1], _passUpdate.OldPassword, _passUpdate.NewPassword)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully updated password"))
}

func (restAuthHandler RestAuthHandlerAdapter) GetUpdateProfile(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.Header.Get("Authorization"), " ")
	if token == nil || len(token) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not Authorized"))
		return
	}

	var _user auth_domain.User

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&_user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	defer r.Body.Close()

	err = restAuthHandler.authService.UpdateProfile(token[1], _user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("successfully updated profile"))
}

func (handler RestAuthHandlerAdapter) GetStudents(w http.ResponseWriter, r *http.Request) {
	token := strings.Split(r.Header.Get("Authorization"), " ")
	if token == nil || len(token) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not Authorized"))
		return
	}

	students, err := handler.authService.GetStudents(token[1])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	sendMasrshaledResponse(students, w)

}
func (handler RestAuthHandlerAdapter) GetUpdateStudentPayment(w http.ResponseWriter, r *http.Request) {
	var userId string
	var status bool = false
	token := strings.Split(r.Header.Get("Authorization"), " ")
	if token == nil || len(token) < 2 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not Authorized"))
		return
	}

	userId = r.URL.Query().Get("id")
	if userId == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Not valid user id"))
		return
	}
	if r.URL.Query().Get("status") != "" {
		status = true
	}

	student, err := handler.authService.UpdateStudentPayment(token[1], userId, status)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	sendMasrshaledResponse(student, w)

}
