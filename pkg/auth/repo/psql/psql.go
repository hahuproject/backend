package auth_psql_repo

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	auth_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/domain"
	authRepo "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/repo"
	authError "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/utils"

	"golang.org/x/crypto/bcrypt"
)

type PsqlAuthRepositoryAdapter struct {
	log *log.Logger
	db  *sql.DB
}

func NewPsqlAuthRepositoryAdapter(log *log.Logger, db *sql.DB) (authRepo.AuthRepositoryPort, error) {
	err := db.Ping()
	if err != nil {
		log.Println(err)
		log.Println("DB Ping failed")
		return &PsqlAuthRepositoryAdapter{}, err
	}

	return &PsqlAuthRepositoryAdapter{
		log: log,
		db:  db,
	}, nil
}

func StoreAddress(db *sql.DB, userAddress auth_domain.Address) (auth_domain.Address, error) {
	var address = auth_domain.Address{}

	log.Println(userAddress)

	err := db.QueryRow("INSERT INTO public.addresses (country, region, city, subcity, woreda, house_no) VALUES ($1,$2,$3,$4,$5,$6) RETURNING *", userAddress.Country, userAddress.Region, userAddress.City, userAddress.SubCity, userAddress.Woreda, userAddress.HouseNo).Scan(&address.ID, &address.Country, &address.Region, &address.City, &address.SubCity, &address.Woreda, &address.HouseNo)

	if err != nil {
		return auth_domain.Address{}, errors.New("failed to save data (Address)")
	}
	return address, nil
}

func (psqlAdapter *PsqlAuthRepositoryAdapter) StoreAdmin(user auth_domain.Admin) (auth_domain.Admin, error) {

	var addedUser auth_domain.Admin

	//Store Address
	address, err := StoreAddress(psqlAdapter.db, user.User.Address)
	if err != nil {
		return addedUser, err
	}

	//Store User
	rows, err := psqlAdapter.db.Query("SELECT * FROM public.user_types WHERE type = 'ADMIN'")
	if err != nil {
		return addedUser, err
	}
	defer rows.Close()

	type Admin struct {
		ID     int
		Type   string
		UserId string
	}
	var admins []Admin

	for rows.Next() {
		var admin Admin
		err := rows.Scan(&admin.ID, &admin.UserId, &admin.Type)
		if err != nil {
			// ////psqlAdapter.log.Println(err)
			return addedUser, err
		}
		// //psqlAdapter.log.Println(admin)

		admins = append(admins, admin)
	}
	if err = rows.Err(); err != nil {
		return addedUser, err
	}

	//ID generation
	//school, employment, type/4 digit/ year

	//psqlAdapter.log.Println(admins)
	count := fmt.Sprintf("%04d", len(admins)+1)
	newUserId := "TFA" + count + strconv.Itoa(time.Now().Year())
	//psqlAdapter.log.Println(newUserId)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.User.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	err = psqlAdapter.db.QueryRow("INSERT INTO public.users (user_id, first_name, last_name, email, phone, profile_pic, username, password, address) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING user_id, first_name, last_name, email, phone, username, address, profile_pic", newUserId, user.User.FirstName, user.User.LastName, user.User.Email, user.User.Phone, user.User.ProfilePic, user.User.Username, hashedPassword, address.ID).Scan(&addedUser.ID, &addedUser.User.FirstName, &addedUser.User.LastName, &addedUser.User.Email, &addedUser.User.Phone, &addedUser.User.Username, &addedUser.User.Address.ID, &addedUser.User.ProfilePic)

	// check for err sql.ErrNoRows
	// check for err sql.ErrNoRows
	if err != nil {
		psqlAdapter.db.QueryRow("DELETE FROM public.addresses WHERE address_id = $1", address.ID)
		////psqlAdapter.log.Println(err)
		if err.Error() == `pq: duplicate key value violates unique constraint "username"` {
			//psqlAdapter.log.Println("true")
			return addedUser, authError.ErrUsernameTaken
		}
		if err.Error() == "pq: duplicate key value violates unique constraint \"email\"" {
			return addedUser, authError.ErrEmailTaken
		}
		return addedUser, errors.New("failed to save data (User)" + err.Error())
	}

	//Store User Types
	psqlAdapter.db.QueryRow("INSERT INTO public.user_types (user_id, type) VALUES ($1, 'ADMIN')", addedUser.ID)

	//Store Specific User
	psqlAdapter.db.QueryRow("INSERT INTO public.admins (user_id) VALUES ($1)", addedUser.ID).Scan()

	return addedUser, nil
}

func (psqlAdapter *PsqlAuthRepositoryAdapter) StoreRegistryOfficer(user auth_domain.RegistryOfficer) (auth_domain.RegistryOfficer, error) {

	var addedUser auth_domain.RegistryOfficer

	//Store Address
	address, err := StoreAddress(psqlAdapter.db, user.User.Address)
	if err != nil {
		return addedUser, err
	}

	//Store User
	rows, err := psqlAdapter.db.Query("SELECT * FROM public.user_types WHERE type = 'REGISTRY_OFFICER'")
	if err != nil {
		return addedUser, err
	}
	defer rows.Close()

	type RegOfficer struct {
		ID     int
		Type   string
		UserId string
	}
	var regOfficers []RegOfficer

	for rows.Next() {
		var regOfficer RegOfficer
		err := rows.Scan(&regOfficer.ID, &regOfficer.UserId, &regOfficer.Type)
		if err != nil {
			return addedUser, err
		}

		regOfficers = append(regOfficers, regOfficer)
	}
	if err = rows.Err(); err != nil {
		return addedUser, err
	}

	//ID generation
	//school, employment, type/4 digit/ year

	//psqlAdapter.log.Println(regOfficers)
	count := fmt.Sprintf("%04d", len(regOfficers)+1)
	newUserId := "TFRO" + count + strconv.Itoa(time.Now().Year())

	// var userID string
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.User.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	err = psqlAdapter.db.QueryRow("INSERT INTO public.users (user_id, first_name, last_name, email, phone, profile_pic, username, password, address) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING user_id, first_name, last_name, email, phone, username, address, profile_pic", newUserId, user.User.FirstName, user.User.LastName, user.User.Email, user.User.Phone, user.User.ProfilePic, user.User.Username, hashedPassword, address.ID).Scan(&addedUser.ID, &addedUser.User.FirstName, &addedUser.User.LastName, &addedUser.User.Email, &addedUser.User.Phone, &addedUser.User.Username, &addedUser.User.Address.ID, &addedUser.User.ProfilePic)

	// check for err sql.ErrNoRows
	if err != nil {
		psqlAdapter.db.QueryRow("DELETE FROM public.addresses WHERE address_id = $1", address.ID)
		////psqlAdapter.log.Println(err)
		if err.Error() == `pq: duplicate key value violates unique constraint "username"` {
			//psqlAdapter.log.Println("true")
			return addedUser, authError.ErrUsernameTaken
		}
		if err.Error() == "pq: duplicate key value violates unique constraint \"email\"" {
			return addedUser, authError.ErrEmailTaken
		}
		return addedUser, errors.New("failed to save data (User)" + err.Error())
	}

	//Store User Types
	psqlAdapter.db.QueryRow("INSERT INTO public.user_types (user_id, type) VALUES ($1, 'REGISTRY_OFFICER')", addedUser.ID)

	//Store Specific User
	psqlAdapter.db.QueryRow("INSERT INTO public.registry_officers (user_id) VALUES ($1)", addedUser.ID).Scan()

	return addedUser, nil
}

func (psqlAdapter *PsqlAuthRepositoryAdapter) StoreSubRegistryOfficer(user auth_domain.SubRegistryOfficer) (auth_domain.SubRegistryOfficer, error) {

	var addedUser auth_domain.SubRegistryOfficer

	//Store Address
	address, err := StoreAddress(psqlAdapter.db, user.User.Address)
	if err != nil {
		return addedUser, err
	}

	//Store User
	rows, err := psqlAdapter.db.Query("SELECT * FROM public.user_types WHERE type = 'SUB_REGISTRY_OFFICER'")
	if err != nil {
		return addedUser, err
	}
	defer rows.Close()

	type RegOfficer struct {
		ID     int
		Type   string
		UserId string
	}
	var regOfficers []RegOfficer

	for rows.Next() {
		var regOfficer RegOfficer
		err := rows.Scan(&regOfficer.ID, &regOfficer.UserId, &regOfficer.Type)
		if err != nil {
			return addedUser, err
		}

		regOfficers = append(regOfficers, regOfficer)
	}
	if err = rows.Err(); err != nil {
		return addedUser, err
	}

	//ID generation
	//school, employment, type/4 digit/ year

	//psqlAdapter.log.Println(regOfficers)
	count := fmt.Sprintf("%04d", len(regOfficers)+1)
	newUserId := "TFSRO" + count + strconv.Itoa(time.Now().Year())

	// var userID string
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.User.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	err = psqlAdapter.db.QueryRow("INSERT INTO public.users (user_id, first_name, last_name, email, phone, profile_pic, username, password, address) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING user_id, first_name, last_name, email, phone, username, address, profile_pic", newUserId, user.User.FirstName, user.User.LastName, user.User.Email, user.User.Phone, user.User.ProfilePic, user.User.Username, hashedPassword, address.ID).Scan(&addedUser.ID, &addedUser.User.FirstName, &addedUser.User.LastName, &addedUser.User.Email, &addedUser.User.Phone, &addedUser.User.Username, &addedUser.User.Address.ID, &addedUser.User.ProfilePic)

	// check for err sql.ErrNoRows
	if err != nil {
		psqlAdapter.db.QueryRow("DELETE FROM public.addresses WHERE address_id = $1", address.ID)
		////psqlAdapter.log.Println(err)
		if err.Error() == `pq: duplicate key value violates unique constraint "username"` {
			//psqlAdapter.log.Println("true")
			return addedUser, authError.ErrUsernameTaken
		}
		if err.Error() == "pq: duplicate key value violates unique constraint \"email\"" {
			return addedUser, authError.ErrEmailTaken
		}
		return addedUser, errors.New("failed to save data (User)" + err.Error())
	}

	//Store User Types
	psqlAdapter.db.QueryRow("INSERT INTO public.user_types (user_id, type) VALUES ($1, 'REGISTRY_OFFICER')", addedUser.ID)

	//Store Specific User
	psqlAdapter.db.QueryRow("INSERT INTO public.sub_registry_officers (user_id, department_id) VALUES ($1, $2)", addedUser.ID, user.Department.ID).Scan()

	return addedUser, nil
}

func (psqlAdapter *PsqlAuthRepositoryAdapter) StoreDepartmentHead(user auth_domain.DepartmentHead) (auth_domain.DepartmentHead, error) {
	var addedUser auth_domain.DepartmentHead

	//Store Address
	address, err := StoreAddress(psqlAdapter.db, user.User.Address)
	if err != nil {
		return addedUser, err
	}

	//Store User
	rows, err := psqlAdapter.db.Query("SELECT * FROM public.user_types WHERE type = 'DEPARTMENT_HEAD'")
	if err != nil {
		return addedUser, err
	}
	defer rows.Close()

	type DepHead struct {
		ID     int
		Type   string
		UserId string
	}
	var depHeads []DepHead

	for rows.Next() {
		var depHead DepHead
		err := rows.Scan(&depHead.ID, &depHead.UserId, &depHead.Type)
		if err != nil {
			return addedUser, err
		}

		depHeads = append(depHeads, depHead)
	}
	if err = rows.Err(); err != nil {
		return addedUser, err
	}

	//ID generation
	//school, employment, type/4 digit/ year

	//psqlAdapter.log.Println(depHeads)
	count := fmt.Sprintf("%04d", len(depHeads)+1)
	newUserId := "TFDH" + count + strconv.Itoa(time.Now().Year())

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.User.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	err = psqlAdapter.db.QueryRow("INSERT INTO public.users (user_id, first_name, last_name, email, phone, profile_pic, username, password, address) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING user_id, first_name, last_name, email, phone, username, address, profile_pic", newUserId, user.User.FirstName, user.User.LastName, user.User.Email, user.User.Phone, user.User.ProfilePic, user.User.Username, hashedPassword, address.ID).Scan(&addedUser.ID, &addedUser.User.FirstName, &addedUser.User.LastName, &addedUser.User.Email, &addedUser.User.Phone, &addedUser.User.Username, &addedUser.User.Address.ID, &addedUser.User.ProfilePic)

	// check for err sql.ErrNoRows
	if err != nil {
		psqlAdapter.db.QueryRow("DELETE FROM public.addresses WHERE address_id = $1", address.ID)
		////psqlAdapter.log.Println(err)
		if err.Error() == `pq: duplicate key value violates unique constraint "username"` {
			//psqlAdapter.log.Println("true")
			return addedUser, authError.ErrUsernameTaken
		}
		if err.Error() == "pq: duplicate key value violates unique constraint \"email\"" {
			return addedUser, authError.ErrEmailTaken
		}
		return addedUser, errors.New("failed to save data (User)" + err.Error())
	}

	//Store User Types
	psqlAdapter.db.QueryRow("INSERT INTO public.user_types (user_id, type) VALUES ($1, 'DEPARTMENT_HEAD')", addedUser.ID)

	//Store Specific User
	psqlAdapter.db.QueryRow("INSERT INTO public.department_heads (user_id) VALUES ($1)", addedUser.ID).Scan()

	return addedUser, nil
}
func (psqlAdapter *PsqlAuthRepositoryAdapter) StoreInstructor(user auth_domain.Instructor) (auth_domain.Instructor, error) {

	var addedUser auth_domain.Instructor

	//Store Address
	address, err := StoreAddress(psqlAdapter.db, user.User.Address)
	if err != nil {
		return addedUser, err
	}

	//Store User
	rows, err := psqlAdapter.db.Query("SELECT * FROM public.user_types WHERE type = 'INSTRUCTOR'")
	if err != nil {
		return addedUser, err
	}
	defer rows.Close()

	type Instructor struct {
		ID     int
		Type   string
		UserId string
	}
	var instructors []Instructor

	for rows.Next() {
		var instructor Instructor
		err := rows.Scan(&instructor.ID, &instructor.UserId, &instructor.Type)
		if err != nil {
			return addedUser, err
		}

		instructors = append(instructors, instructor)
	}
	if err = rows.Err(); err != nil {
		return addedUser, err
	}

	//ID generation
	//school, employment, type/4 digit/ year

	//psqlAdapter.log.Println(instructors)
	count := fmt.Sprintf("%04d", len(instructors)+1)
	var newUserId string = "T" + strings.ToUpper(string(user.EmploymentType[0])) + "I" + count + strconv.Itoa(time.Now().Year())

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.User.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	err = psqlAdapter.db.QueryRow("INSERT INTO public.users (user_id, first_name, last_name, email, phone, profile_pic, username, password, address) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING user_id, first_name, last_name, email, phone, username, address, profile_pic", newUserId, user.User.FirstName, user.User.LastName, user.User.Email, user.User.Phone, user.User.ProfilePic, user.User.Username, hashedPassword, address.ID).Scan(&addedUser.ID, &addedUser.User.FirstName, &addedUser.User.LastName, &addedUser.User.Email, &addedUser.User.Phone, &addedUser.User.Username, &addedUser.User.Address.ID, &addedUser.User.ProfilePic)

	// check for err sql.ErrNoRows
	if err != nil {
		psqlAdapter.db.QueryRow("DELETE FROM public.addresses WHERE address_id = $1", address.ID)
		////psqlAdapter.log.Println(err)
		if err.Error() == `pq: duplicate key value violates unique constraint "username"` {
			//psqlAdapter.log.Println("true")
			return addedUser, authError.ErrUsernameTaken
		}
		if err.Error() == "pq: duplicate key value violates unique constraint \"email\"" {
			return addedUser, authError.ErrEmailTaken
		}
		return addedUser, errors.New("failed to save data (User)" + err.Error())
	}

	//Store User Types
	psqlAdapter.db.QueryRow("INSERT INTO public.user_types (user_id, type) VALUES ($1, 'INSTRUCTOR')", addedUser.ID)

	//Store Specific User
	psqlAdapter.db.QueryRow("INSERT INTO public.instructors (user_id, salary_rate, employment_type) VALUES ($1,$2,$3)", addedUser.ID, 0, user.EmploymentType).Scan()

	return addedUser, nil
}
func (psqlAdapter *PsqlAuthRepositoryAdapter) StoreSupervisor(user auth_domain.Supervisor) (auth_domain.Supervisor, error) {
	var addedUser auth_domain.Supervisor

	//Store Address
	address, err := StoreAddress(psqlAdapter.db, user.User.Address)
	if err != nil {
		return addedUser, err
	}

	//Store User
	rows, err := psqlAdapter.db.Query("SELECT * FROM public.user_types WHERE type = 'SUPERVISOR'")
	if err != nil {
		return addedUser, err
	}
	defer rows.Close()

	type Supervisor struct {
		ID     int
		Type   string
		UserId string
	}
	var supervisors []Supervisor

	for rows.Next() {
		var supervisor Supervisor
		err := rows.Scan(&supervisor.ID, &supervisor.UserId, &supervisor.Type)
		if err != nil {
			return addedUser, err
		}

		supervisors = append(supervisors, supervisor)
	}
	if err = rows.Err(); err != nil {
		return addedUser, err
	}

	//ID generation
	//school, employment, type/4 digit/ year

	//psqlAdapter.log.Println(supervisors)
	count := fmt.Sprintf("%04d", len(supervisors)+1)
	newUserId := "TFS" + count + strconv.Itoa(time.Now().Year())

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.User.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	err = psqlAdapter.db.QueryRow("INSERT INTO public.users (user_id, first_name, last_name, email, phone, profile_pic, username, password, address) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING user_id, first_name, last_name, email, phone, username, address, profile_pic", newUserId, user.User.FirstName, user.User.LastName, user.User.Email, user.User.Phone, user.User.ProfilePic, user.User.Username, hashedPassword, address.ID).Scan(&addedUser.ID, &addedUser.User.FirstName, &addedUser.User.LastName, &addedUser.User.Email, &addedUser.User.Phone, &addedUser.User.Username, &addedUser.User.Address.ID, &addedUser.User.ProfilePic)

	// check for err sql.ErrNoRows
	if err != nil {
		psqlAdapter.db.QueryRow("DELETE FROM public.addresses WHERE address_id = $1", address.ID)
		////psqlAdapter.log.Println(err)
		if err.Error() == `pq: duplicate key value violates unique constraint "username"` {
			//psqlAdapter.log.Println("true")
			return addedUser, authError.ErrUsernameTaken
		}
		if err.Error() == "pq: duplicate key value violates unique constraint \"email\"" {
			return addedUser, authError.ErrEmailTaken
		}
		return addedUser, errors.New("failed to save data (User)" + err.Error())
	}

	//Store User Types
	psqlAdapter.db.QueryRow("INSERT INTO public.user_types (user_id, type) VALUES ($1, 'SUPERVISOR')", addedUser.ID)

	//Store Specific User
	psqlAdapter.db.QueryRow("INSERT INTO public.supervisors (user_id, industry) VALUES ($1,$2)", addedUser.ID, user.Industry).Scan()

	return addedUser, nil
}
func substr(input string, start int, length int) string {
	asRunes := []rune(input)

	if start >= len(asRunes) {
		return ""
	}

	if start+length > len(asRunes) {
		length = len(asRunes) - start
	}

	return string(asRunes[start : start+length])
}

func (psqlAdapter *PsqlAuthRepositoryAdapter) StoreStudent(verified bool, user auth_domain.Student) (auth_domain.Student, error) {
	var addedUser auth_domain.Student

	//Store Address
	address, err := StoreAddress(psqlAdapter.db, user.User.Address)
	if err != nil {
		return addedUser, err
	}
	emergencyContactAddress, err := StoreAddress(psqlAdapter.db, user.EmergencyContactAddress)
	if err != nil {
		return addedUser, err
	}

	//school dpt type year
	//Store User
	rows, err := psqlAdapter.db.Query("SELECT user_type_id, type FROM public.user_types WHERE type = 'STUDENT'")
	if err != nil {
		return addedUser, err
	}
	defer rows.Close()

	type Student struct {
		ID   int
		Type string
	}
	var students []Student

	for rows.Next() {
		var student Student
		err := rows.Scan(&student.ID, &student.Type)
		if err != nil {
			return addedUser, err
		}

		students = append(students, student)
	}
	if err = rows.Err(); err != nil {
		return addedUser, err
	}

	//ID generation
	//school, employment, type/4 digit/ year
	//TRSER/0001/22

	//psqlAdapter.log.Println(len(students))
	// psqlAdapter.log.Println("user.Class.Department.Name")
	// psqlAdapter.log.Panicln(user.Class.Department.Name)
	count := fmt.Sprintf("%04d", len(students)+1)
	var newUserId string = "TR" + substr(strings.ToUpper(string(user.Department.Name)), 0, 2) + substr(strings.ToUpper(string(user.Program)), 0, 1) + "/" + count + "/" + substr(strconv.Itoa(time.Now().Year()), 2, 4)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.User.Password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	err = psqlAdapter.db.QueryRow("INSERT INTO public.users (user_id, first_name, last_name, email, phone, profile_pic, username, password, address, verified) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING user_id, first_name, last_name, email, phone, username, address, profile_pic", newUserId, user.User.FirstName, user.User.LastName, user.User.Email, user.User.Phone, user.User.ProfilePic, user.User.Username, hashedPassword, address.ID, verified).Scan(&addedUser.ID, &addedUser.User.FirstName, &addedUser.User.LastName, &addedUser.User.Email, &addedUser.User.Phone, &addedUser.User.Username, &addedUser.User.Address.ID, &addedUser.User.ProfilePic)

	// check for err sql.ErrNoRows
	if err != nil {
		psqlAdapter.db.QueryRow("DELETE FROM public.addresses WHERE address_id = $1", address.ID)
		////psqlAdapter.log.Println(err)
		if err.Error() == `pq: duplicate key value violates unique constraint "username"` {
			//psqlAdapter.log.Println("true")
			return addedUser, authError.ErrUsernameTaken
		}
		if err.Error() == "pq: duplicate key value violates unique constraint \"email\"" {
			return addedUser, authError.ErrEmailTaken
		}
		return addedUser, errors.New("failed to save data (User)" + err.Error())
	}

	//Store User Types
	psqlAdapter.db.QueryRow("INSERT INTO public.user_types (user_id, type) VALUES ($1, 'STUDENT')", addedUser.ID)

	//Store Specific User
	err = psqlAdapter.db.QueryRow(`INSERT INTO public.students 
	(user_id, middle_name, gender, disability, birth_date, birth_place, 
		previous_school, program, matric_result, average_mark_for_highschool, 
		emergency_contact_name, emergency_contact_phone, emergency_contact_relation, emergency_contact_address, 
		department_id, stream_id) 
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15, $16)`,
		addedUser.ID, user.MiddleName, user.Gender, user.Disablility, user.BirthDate, user.BirthPlace,
		user.PreviousSchool, user.Program, user.MatricResult, user.AveragemarkForHighSchool,
		user.EmergencyContactName, user.EmergencyContactPhone, user.EmergencyContactRelation, emergencyContactAddress.ID,
		sql.NullString{Valid: user.Department.ID != "", String: user.Department.ID}, sql.NullString{Valid: user.Stream.ID != "", String: user.Stream.ID}).Scan()

	if err != nil {
		psqlAdapter.log.Println("err in save students")
		psqlAdapter.log.Println(err)
		//psqlAdapter.log.Println(user.MatricResult)
		//psqlAdapter.log.Println(user.AveragemarkForHighSchool)
	}

	return addedUser, nil
}
func (psqlAdapter *PsqlAuthRepositoryAdapter) FindUserByUsername(username string) (auth_domain.User, error) {
	user := auth_domain.User{}

	err := psqlAdapter.db.QueryRow(`
	SELECT 
	user_id, first_name, last_name, email, phone, username, password, verified, banned 
	FROM public.users WHERE username = $1`, username).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Username, &user.Password, &user.Verified, &user.Banned)
	if err == sql.ErrNoRows {
		return auth_domain.User{}, authError.ErrUserNotFound
	} else if err != nil {
		return auth_domain.User{}, err
	}
	return user, nil
}

func (psqlAdapter *PsqlAuthRepositoryAdapter) FindUserById(id string) (auth_domain.User, error) {
	user := auth_domain.User{}

	err := psqlAdapter.db.QueryRow(`SELECT 
	public.users.user_id, first_name, last_name, email, phone, username, password, profile_pic, verified, type,
	public.addresses.address_id,public.addresses.country, public.addresses.region, public.addresses.city, public.addresses.subcity,public.addresses.woreda, public.addresses.house_no
	FROM public.users 
	INNER JOIN public.user_types 
	ON public.users.user_id = public.user_types.user_id 
	INNER JOIN public.addresses
	ON public.users.address = public.addresses.address_id
	WHERE public.users.user_id = $1`, id).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Username, &user.Password, &user.ProfilePic, &user.Verified, &user.Type,
		&user.Address.ID, &user.Address.Country, &user.Address.Region, &user.Address.City, &user.Address.SubCity, &user.Address.Woreda, &user.Address.HouseNo)
	if err == sql.ErrNoRows {
		return auth_domain.User{}, authError.ErrUserNotFound
	} else if err != nil {
		return auth_domain.User{}, err
	}
	return user, nil
}

func (psqlAdapter *PsqlAuthRepositoryAdapter) ApproveUser(id string) error {
	_, err := psqlAdapter.db.Query("UPDATE public.users SET verified = 'true' WHERE user_id = $1", id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return authError.ErrUserNotFound
		}
		return err
	}
	return nil
}
func (psqlAdapter *PsqlAuthRepositoryAdapter) BanUser(id string) error {
	_, err := psqlAdapter.db.Exec(`UPDATE public.users SET banned = $1 WHERE user_id = $2`, true, id)
	return err
}
func (psqlAdapter *PsqlAuthRepositoryAdapter) UnBanUser(id string) error {
	_, err := psqlAdapter.db.Exec(`UPDATE public.users SET banned = $1 WHERE user_id = $2`, false, id)
	return err
}
func (psqlAdapter *PsqlAuthRepositoryAdapter) GetUserRoleByID(id string) (string, error) {
	var role string
	err := psqlAdapter.db.QueryRow("SELECT type FROM public.user_types WHERE user_id = $1", id).Scan(&role)
	if err != nil {
		////psqlAdapter.log.Println(err)
		return "", err
	}
	//psqlAdapter.log.Println(role)

	return role, nil
}

func (psqlAdapter *PsqlAuthRepositoryAdapter) FindAllUsers() ([]auth_domain.User, error) {

	var users []auth_domain.User

	rows, err := psqlAdapter.db.Query(`
	SELECT 
	public.users.user_id, first_name, last_name, email, phone, username, profile_pic, address, verified, type, created_at, banned,
	public.addresses.country, public.addresses.region, public.addresses.city, public.addresses.subcity, public.addresses.woreda, public.addresses.house_no 
	FROM public.users 
	INNER JOIN public.user_types ON public.users.user_id = public.user_types.user_id
	INNER JOIN public.addresses ON public.users.address = public.addresses.address_id
	`)
	if err != nil {
		psqlAdapter.log.Println(err)
		return users, authError.ErrNoUsersFound
	}

	for rows.Next() {
		var _user auth_domain.User
		rows.Scan(&_user.ID, &_user.FirstName, &_user.LastName, &_user.Email, &_user.Phone, &_user.Username, &_user.ProfilePic, &_user.Address.ID, &_user.Verified, &_user.Type, &_user.CreatedAt, &_user.Banned,
			&_user.Address.Country, &_user.Address.Region, &_user.Address.City, &_user.Address.SubCity, &_user.Address.Woreda, &_user.Address.HouseNo)

		users = append(users, _user)
	}

	return users, nil
}

func populateStudent(db *sql.DB, student *auth_domain.Student) {
	student.Courses = make([]auth_domain.Course, 0)

	rows, err := db.Query(`
	SELECT public.students_courses.course_id, public.courses.name, public.courses.credit_hr, public.courses.color
	FROM public.students_courses
	INNER JOIN public.courses ON public.students_courses.course_id = public.courses.course_id
	WHERE public.students_courses.student_id = $1
	`, student.ID)
	if err == nil {
		defer rows.Close()

		for rows.Next() {
			var _course auth_domain.Course
			rows.Scan(&_course.ID, &_course.Name, &_course.CreditHr, &_course.Color)

			student.Courses = append(student.Courses, _course)
		}
	}
}

func (repo PsqlAuthRepositoryAdapter) FindStudent(id string) (auth_domain.Student, error) {
	var _student auth_domain.Student
	// public.students.emergency_contact_name, public.students.emergency_contact_phone, public.students.emergency_contact_relation, public.students.emergency_contact_address,

	err := repo.db.QueryRow(`
	SELECT
		public.students.user_id, public.students.middle_name, public.students.gender, public.students.disability, public.students.birth_date, public.students.birth_place, paid,
		public.students.previous_school, public.students.program, public.students.matric_result, public.students.average_mark_for_highschool, 
		public.users.user_id, first_name, last_name, email, phone, username, profile_pic, verified, type,
		public.addresses.address_id,public.addresses.country, public.addresses.region, public.addresses.city, public.addresses.subcity,public.addresses.woreda, public.addresses.house_no,
		public.students.department_id, public.departments.name
	FROM public.students
	INNER JOIN public.users ON public.users.user_id = public.students.user_id
	INNER JOIN public.user_types ON public.users.user_id = public.user_types.user_id
	INNER JOIN public.addresses ON public.users.address = public.addresses.address_id
	LEFT JOIN public.departments ON public.students.department_id = public.departments.department_id
	WHERE public.students.user_id = $1
	`, id).Scan(
		&_student.ID, &_student.MiddleName, &_student.Gender, &_student.Disablility, &_student.BirthDate, &_student.BirthPlace, &_student.Paid,
		&_student.PreviousSchool, &_student.Program, &_student.MatricResult, &_student.AveragemarkForHighSchool,
		&_student.User.ID, &_student.User.FirstName, &_student.User.LastName, &_student.User.Email, &_student.User.Phone, &_student.User.Username, &_student.User.ProfilePic, &_student.User.Verified, &_student.User.Type,
		&_student.User.Address.ID, &_student.User.Address.Country, &_student.User.Address.Region, &_student.User.Address.City, &_student.User.Address.SubCity, &_student.User.Address.Woreda, &_student.User.Address.HouseNo,
		&_student.Department.ID, &_student.Department.Name,
	)
	if err != nil {
		return _student, err
	}

	populateStudent(repo.db, &_student)

	return _student, nil
}
func (repo PsqlAuthRepositoryAdapter) FindInstructor(id string) (auth_domain.Instructor, error) {
	var _instructor auth_domain.Instructor

	err := repo.db.QueryRow(`
	SELECT
		public.instructors.user_id, public.instructors.salary_rate, public.instructors.employment_type,
		public.users.user_id, first_name, last_name, email, phone, username, profile_pic, verified, type,
		public.addresses.address_id,public.addresses.country, public.addresses.region, public.addresses.city, public.addresses.subcity,public.addresses.woreda, public.addresses.house_no
	FROM public.instructors
	INNER JOIN public.users ON public.users.user_id = public.instructors.user_id
	INNER JOIN public.user_types ON public.users.user_id = public.user_types.user_id
	INNER JOIN public.addresses ON public.users.address = public.addresses.address_id
	WHERE public.instructors.user_id = $1
	`, id).Scan(
		&_instructor.ID, &_instructor.SalaryRate, &_instructor.EmploymentType,
		&_instructor.User.ID, &_instructor.User.FirstName, &_instructor.User.LastName, &_instructor.User.Email, &_instructor.User.Phone, &_instructor.User.Username, &_instructor.User.ProfilePic, &_instructor.User.Verified, &_instructor.User.Type,
		&_instructor.User.Address.ID, &_instructor.User.Address.Country, &_instructor.User.Address.Region, &_instructor.User.Address.City, &_instructor.User.Address.SubCity, &_instructor.User.Address.Woreda, &_instructor.User.Address.HouseNo,
	)
	if err != nil {
		return _instructor, err
	}

	return _instructor, nil
}
func (repo PsqlAuthRepositoryAdapter) FindSupervisor(id string) (auth_domain.Supervisor, error) {
	var _supervisor auth_domain.Supervisor

	err := repo.db.QueryRow(`
	SELECT
		public.supervisors.user_id, public.supervisors.industry,
		public.users.user_id, first_name, last_name, email, phone, username, profile_pic, verified, type,
		public.addresses.address_id,public.addresses.country, public.addresses.region, public.addresses.city, public.addresses.subcity,public.addresses.woreda, public.addresses.house_no
	FROM public.supervisors
	INNER JOIN public.users ON public.users.user_id = public.supervisors.user_id
	INNER JOIN public.user_types ON public.users.user_id = public.user_types.user_id
	INNER JOIN public.addresses ON public.users.address = public.addresses.address_id
	WHERE public.supervisors.user_id = $1
	`, id).Scan(
		&_supervisor.ID, &_supervisor.Industry,
		&_supervisor.User.ID, &_supervisor.User.FirstName, &_supervisor.User.LastName, &_supervisor.User.Email, &_supervisor.User.Phone, &_supervisor.User.Username, &_supervisor.User.ProfilePic, &_supervisor.User.Verified, &_supervisor.User.Type,
		&_supervisor.User.Address.ID, &_supervisor.User.Address.Country, &_supervisor.User.Address.Region, &_supervisor.User.Address.City, &_supervisor.User.Address.SubCity, &_supervisor.User.Address.Woreda, &_supervisor.User.Address.HouseNo,
	)
	if err != nil {
		return _supervisor, err
	}

	return _supervisor, nil
}

func (repo PsqlAuthRepositoryAdapter) ChangePassword(userId, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = repo.db.Exec(`UPDATE public.users SET password = $1 WHERE user_id = $2`, hashedPassword, userId)
	return err
}

func (repo PsqlAuthRepositoryAdapter) UpdateProfile(user auth_domain.User) error {
	_, err := repo.db.Exec(`
	UPDATE public.users SET 
	first_name = $1,
	last_name = $2,
	email = $3,
	phone = $4
	WHERE user_id = $5`, user.FirstName, user.LastName, user.Email, user.Phone, user.ID)
	return err
}

func (repo PsqlAuthRepositoryAdapter) FindStudents() ([]auth_domain.Student, error) {
	var students []auth_domain.Student = make([]auth_domain.Student, 0)
	rows, err := repo.db.Query(`SELECT user_id FROM public.students`)
	if err != nil {
		return students, err
	}

	for rows.Next() {
		var student auth_domain.Student
		rows.Scan(&student.ID)

		student, err = repo.FindStudent(student.ID)
		if err != nil {
			repo.log.Println(err)
		}

		students = append(students, student)
	}

	return students, err
}
func (repo PsqlAuthRepositoryAdapter) UpdateStudentPaymentSatus(userId string, status bool) (auth_domain.Student, error) {
	_, err := repo.db.Query(`UPDATE public.students SET paid = $1 WHERE user_id = $2`, status, userId)
	if err != nil {
		return auth_domain.Student{}, err
	}

	return repo.FindStudent(userId)
}
