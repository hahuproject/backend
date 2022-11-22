package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	//Auth Imports

	auth_rest_handler "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/handler/rest"
	auth_psql_repo "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/repo/psql"
	auth_service "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/service"
	certificate_rest_handler "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/certificate/handler/rest"
	certificate_service "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/certificate/service"

	//Class Imports

	class_rest_handler "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/handler/rest"
	class_psql_repo "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/repo/psql"
	class_service "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/class/service"

	//Announcement imports
	annoucement_rest_handler "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/announcement/handler/rest"
	annoucement_psql_repo "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/announcement/repo/psql"
	announcement_service "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/announcement/service"

	//Message Imports

	message_rest_handler "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/message/handler/rest"
	message_psql_repo "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/message/repo/psql"
	message_service "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/message/service"

	//Assignment Imports
	assignment_rest_handler "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/assignment/handler/rest"
	assignment_psql_repo "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/assignment/repo/psql"
	assignment_service "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/assignment/service"

	//Attendance Imports
	attendance_rest_handler "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/attendance/handler/rest"
	attendance_psql_repo "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/attendance/repo/psql"
	attendance_service "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/attendance/service"

	//Grade Imports
	grade_rest_handler "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/grade/handler/rest"
	grade_psql_repo "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/grade/repo/psql"
	grade_service "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/grade/service"

	//Library Imports

	library_rest_handler "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/library/handler/rest"
	library_psql_repo "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/library/repo/psql"
	library_service "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/library/service"
)

func main() {
	log := log.New(os.Stdout, "HAHU_SMS_SERVER ", log.LstdFlags|log.Lshortfile)

	//////////////////////////Setup Database////////////////////////////////////////

	const (
		// host     = "localhost"
		// port     = 5432
		// user     = "postgres"
		// password = "root"
		// // password = "pass@123"
		// dbname = "ha_f"
		host     = "hahudatabase.cluster-ciuxzlrsiqkq.eu-west-2.rds.amazonaws.com"
		port     = 5432
		user     = "postgres"
		password = "passwordhahu"
		dbname   = "hahudatabase"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	//require
	//disable
	dbase, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Db connection failure : %v", err)
	} else {
		log.Println("Successfully connected to DB")
	}

	dbase.SetMaxOpenConns(12)
	dbase.SetMaxIdleConns(12)

	defer dbase.Close()

	////////////////////////////////////////////////////////////////////////////////

	/////////////////////////// AUTH SERVICE //////////////////////////////////////

	// DB -> Service -> Handler

	authRepo, err := auth_psql_repo.NewPsqlAuthRepositoryAdapter(log, dbase)
	if err != nil {
		os.Exit(0)
	}
	authService := auth_service.NewAuthService(log, authRepo)

	authHandler := auth_rest_handler.NewRestAuthHandlerAdapter(log, authService)

	/////////////////////////////////////////////////////////////////////////////

	/////////////////////////// CLASS SERVICE //////////////////////////////////////

	// DB -> Service -> Handler

	classRepo, err := class_psql_repo.NewPsqlCourseRepoAdapter(log, dbase)
	if err != nil {
		os.Exit(0)
	}
	classService := class_service.NewClassServiceAdapter(log, classRepo)

	classHandler := class_rest_handler.NewRestClassHandlerAdapter(log, classService)

	/////////////////////////////////////////////////////////////////////////////

	/////////////////////////// LIBRARY SERVICE //////////////////////////////////////

	// DB -> Service -> Handler

	libraryRepo, err := library_psql_repo.NewLibraryPsqlRepoAdapter(log, dbase)
	if err != nil {
		os.Exit(0)
	}
	libraryService := library_service.NewLibraryServiceAdapter(log, libraryRepo)

	libraryHandler := library_rest_handler.NewLibraryRestHandlerAdapter(log, libraryService)

	/////////////////////////////////////////////////////////////////////////////

	/////////////////////////// ANNOUNCEMENT SERVICE //////////////////////////////////////

	// DB -> Service -> Handler

	announcementRepo, err := annoucement_psql_repo.NewPsqlAnnouncementAdapter(log, dbase)
	if err != nil {
		os.Exit(0)
	}
	announcementService := announcement_service.NewAnnouncementServiceAdapter(log, announcementRepo)

	announcementHandler := annoucement_rest_handler.NewRestAnnouncementHanlerAdapter(log, announcementService)

	/////////////////////////////////////////////////////////////////////////////

	/////////////////////////// MESSAGE SERVICE //////////////////////////////////////

	// DB -> Service -> Handler

	messageRepo, err := message_psql_repo.NewMessagePsqlRepoAdapter(log, dbase)
	if err != nil {
		os.Exit(0)
	}
	messageService := message_service.NewMessageServiceAdapter(log, messageRepo)

	messageHandler := message_rest_handler.NewMessageRestHandlerAdapter(log, messageService)

	/////////////////////////////////////////////////////////////////////////////

	/////////////////////////// ASSIGNMENT SERVICE //////////////////////////////////////

	// DB -> Service -> Handler

	assignmentRepo, err := assignment_psql_repo.NewAssignmentPsqlRepoAdapter(log, dbase)
	if err != nil {
		os.Exit(0)
	}
	assignmentService := assignment_service.NewAssignmentServiceAdapter(log, assignmentRepo)

	assignmentHandler := assignment_rest_handler.NewAssignmentRestHandlerAdapter(log, assignmentService)

	/////////////////////////////////////////////////////////////////////////////

	/////////////////////////// ASSIGNMENT SERVICE //////////////////////////////////////

	// DB -> Service -> Handler

	attendanceRepo, err := attendance_psql_repo.NewAttendancePsqlRepoAdapter(log, dbase)
	if err != nil {
		os.Exit(0)
	}
	attendanceService := attendance_service.NewAttendanceServiceAdapter(log, attendanceRepo)

	attendanceHandler := attendance_rest_handler.NewAttendanceRestHandlerAdapter(log, attendanceService)

	/////////////////////////////////////////////////////////////////////////////

	/////////////////////////// GRADE SERVICE //////////////////////////////////////

	// DB -> Service -> Handler

	gradeRepo, err := grade_psql_repo.NewGradePsqlRepoAdapter(log, dbase)
	if err != nil {
		os.Exit(0)
	}
	gradeService := grade_service.NewGradeServiceAdapter(log, gradeRepo)

	gradeHandler := grade_rest_handler.NewGradeRestHandlerAdapter(log, gradeService)

	/////////////////////////////////////////////////////////////////////////////

	/////////////////////////// CERTIFICATE SERVICE //////////////////////////////////////

	// DB -> Service -> Handler

	certificateService := certificate_service.NewCertificateServiceAdapter(log, gradeRepo)

	certificateHandler := certificate_rest_handler.NewRestCertificateHandlerAdapter(log, certificateService)

	/////////////////////////////////////////////////////////////////////////////

	////////////////////////// Setup Server /////////////////////////////////////

	sh := http.NewServeMux()

	sh.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch strings.Split(r.URL.Path, "/")[1] {
		case "":
			{
				http.Redirect(w, r, "https://documenter.getpostman.com/view/9284748/UVRHhNoP#350f9b46-35dd-4e9b-b51c-a5cc2898c461", http.StatusTemporaryRedirect)
				return
			}
		case "uploads":
			{
				// log.Println("uploads")
				// log.Println(r.URL.Path)
				// file, _ := ioutil.ReadFile("../uploads/sample.txt")
				// w.Write(file)
				http.ServeFile(w, r, "./"+r.URL.Path)
				// w.Write([]byte("olla files"))
				return
			}
		case "favicon.ico":
			{
				w.Write([]byte("favicon.ico served"))
				return
			}
		}
	})

	////////////////////////////// STATIC ROUTER ///////////////////////////////////

	// sh.HandleFunc("/uploads/*", func(w http.ResponseWriter, r *http.Request) {

	// })

	//////////////////////////////////////////////////////////////////////////////

	sh.HandleFunc("/auth/register/admin", authHandler.GetRegisterAdmin)
	sh.HandleFunc("/auth/register/registry-officer", authHandler.GetRegisterRegistryOfficer)
	sh.HandleFunc("/auth/register/sub-registry-officer", authHandler.GetRegisterSubRegistryOfficer)
	sh.HandleFunc("/auth/register/department-head", authHandler.GetRegisterDepartmentHead)
	sh.HandleFunc("/auth/register/instructor", authHandler.GetRegisterInstructor)
	sh.HandleFunc("/auth/register/supervisor", authHandler.GetRegisterSupervisor)
	sh.HandleFunc("/auth/register/student", authHandler.GetRegisterStudent)
	sh.HandleFunc("/auth/login", authHandler.GetLogin)
	sh.HandleFunc("/auth/approve", authHandler.GetApproveUser)
	sh.HandleFunc("/auth/ban", authHandler.GetBanUser)
	sh.HandleFunc("/auth/unban", authHandler.GetUnBanUser)
	sh.HandleFunc("/auth/me", authHandler.GetMe)
	sh.HandleFunc("/auth/users", authHandler.GetUsers)
	sh.HandleFunc("/auth/students", authHandler.GetStudents)
	sh.HandleFunc("/auth/students/update-payment", authHandler.GetUpdateStudentPayment)
	sh.HandleFunc("/auth/change-password", authHandler.GetChangePassword)
	sh.HandleFunc("/auth/update", authHandler.GetUpdateProfile)

	////////////////////////////// CLASS ROUTER ///////////////////////////////////

	//Course
	sh.HandleFunc("/courses", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				classHandler.GetAddCourse(w, r)
				return
			}
		case http.MethodGet:
			{
				if len(r.URL.Query()["id"]) > 0 {
					classHandler.GetCourse(w, r)
					return
				}

				classHandler.GetCourses(w, r)
				return

			}
		case http.MethodPatch:
			{
				classHandler.GetUpdateCourse(w, r)
				return
			}
		case http.MethodDelete:
			{
				classHandler.GetDeleteCourse(w, r)
				return
			}
		}
	})

	sh.HandleFunc("/courses/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				classHandler.GetAddUserToCourse(w, r)
				return
			}
		case http.MethodPatch:
			{
				classHandler.GetRemoveUserFromCourse(w, r)
				return
			}
		}
	})

	//Department
	sh.HandleFunc("/departments", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				classHandler.GetAddDepartment(w, r)
				return
			}
		case http.MethodGet:
			{
				if len(r.URL.Query()["id"]) > 0 {
					classHandler.GetDepartment(w, r)
					return
				}
				classHandler.GetDepartments(w, r)
				return
			}
		case http.MethodPatch:
			{
				classHandler.GetUpdateDepartment(w, r)
				return
			}
		case http.MethodDelete:
			{
				classHandler.GetDeleteDepartment(w, r)
				return
			}
		}
	})

	//Class
	sh.HandleFunc("/classes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				classHandler.GetAddClass(w, r)
				return
			}
		case http.MethodGet:
			{
				if len(r.URL.Query()["id"]) > 0 {
					classHandler.GetClass(w, r)
					return
				}
				classHandler.GetClasses(w, r)
				return
			}
		case http.MethodPatch:
			{
				classHandler.GetUpdateClass(w, r)
				return
			}
		case http.MethodDelete:
			{
				classHandler.GetDeleteClass(w, r)
				return
			}
		}
	})

	//Class
	sh.HandleFunc("/classes/master-sheet", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				if len(r.URL.Query()["id"]) > 0 {
					classHandler.GetClass(w, r)
					return
				}
				classHandler.GetClassMasterSheet(w, r)
				return
			}
		}
	})

	//Stream
	sh.HandleFunc("/streams", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				classHandler.GetAddStream(w, r)
				return
			}
		// case http.MethodGet:
		// 	{
		// 		if len(r.URL.Query()["id"]) > 0 {
		// 			classHandler.GetClass(w, r)
		// 			return
		// 		}
		// 		classHandler.GetClasses(w, r)
		// 		return
		// 	}
		case http.MethodPatch:
			{
				classHandler.GetUpdateStream(w, r)
				return
			}
		case http.MethodDelete:
			{
				classHandler.GetDeleteStream(w, r)
				return
			}
		}
	})

	//Class-Course
	sh.HandleFunc("/classes/courses", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				classHandler.GetAddClassCourse(w, r)
				return
			}
		case http.MethodPatch:
			{
				classHandler.GetRemoveClassCourse(w, r)
				return
			}
		}
	})

	//Section
	sh.HandleFunc("/sections", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				classHandler.GetAddSection(w, r)
				return
			}
		case http.MethodGet:
			{
				if len(r.URL.Query()["id"]) > 0 {
					classHandler.GetSection(w, r)
					return
				}
				classHandler.GetSections(w, r)
				return
			}
		case http.MethodPatch:
			{
				classHandler.GetUpdateSection(w, r)
				return
			}
		case http.MethodDelete:
			{
				classHandler.GetDeleteSection(w, r)
				return
			}
		}
	})

	//Section-Users
	sh.HandleFunc("/sections/students", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				classHandler.GetAddStudentToSection(w, r)
				return
			}
		case http.MethodPatch:
			{
				classHandler.GetRemoveStudentFromSection(w, r)
				return
			}
		}
	})

	//Session
	sh.HandleFunc("/sessions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				classHandler.GetSessions(w, r)
				return
			}
		case http.MethodPost:
			{
				classHandler.GetAddSession(w, r)
				return
			}
		case http.MethodDelete:
			{
				classHandler.GetDeleteSession(w, r)
				return
			}
		}
	})

	////////////////////////////// ASSIGNMENT ROUTER ///////////////////////////////////

	sh.HandleFunc("/assignments", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				assignmentHandler.GetAddAssignment(w, r)
				return
			}
		case http.MethodGet:
			{
				assignmentHandler.GetAssignments(w, r)
				return
			}
		case http.MethodPatch:
			{
				assignmentHandler.GetUpdateAssignment(w, r)
				return
			}
		case http.MethodDelete:
			{
				assignmentHandler.GetDeleteAssignment(w, r)
				return
			}
		}
	})

	////////////////////////////// GRADE ROUTER ///////////////////////////////////

	sh.HandleFunc("/grades", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				gradeHandler.GetAddGrade(w, r)
				return
			}
		case http.MethodGet:
			{
				gradeHandler.GetGrades(w, r)
				return
			}
		}
	})

	sh.HandleFunc("/grades/request-review", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				gradeHandler.GetRequestGradeReview(w, r)
				return
			}
		}
	})

	sh.HandleFunc("/grades/request-review/approve", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				gradeHandler.GetApproveGradeReview(w, r)
				return
			}
		}
	})

	sh.HandleFunc("/grades/request-review/reject", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				gradeHandler.GetRejectGradeReview(w, r)
				return
			}
		}
	})

	sh.HandleFunc("/grades/submit", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				gradeHandler.GetSubmitGrade(w, r)
				return
			}
		}
	})

	sh.HandleFunc("/grade-labels", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				gradeHandler.GetAddGradeLabel(w, r)
				return
			}
		case http.MethodGet:
			{
				gradeHandler.GetGradeLabels(w, r)
				return
			}
		case http.MethodDelete:
			{
				gradeHandler.GetRemoveGradeLabels(w, r)
				return
			}
		}
	})

	////////////////////////////// ATTENDANCE ROUTER ///////////////////////////////////

	sh.HandleFunc("/attendances", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				attendanceHandler.GetAddAttendance(w, r)
				return
			}
		case http.MethodGet:
			{
				attendanceHandler.GetAttendances(w, r)
				return
			}
		}
	})

	////////////////////////////// LIBRARY ROUTER ///////////////////////////////////

	sh.HandleFunc("/library/books", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				libraryHandler.GetAddBook(w, r)
				return
			}
		case http.MethodGet:
			{
				if len(r.URL.Query()["id"]) > 0 {
					libraryHandler.GetBook(w, r)
				} else {
					libraryHandler.GetBooks(w, r)
				}
				return
			}
		case http.MethodPatch:
			{
				libraryHandler.GetUpdateBook(w, r)
				return
			}
		case http.MethodDelete:
			{
				libraryHandler.GetDeleteBook(w, r)
				return
			}
		}
	})

	sh.HandleFunc("/library/books/rate", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				libraryHandler.GetUpdateBookRating(w, r)
				return
			}
		}
	})

	////////////////////////////// ANNOUNCEMENT ROUTER //////////////////////////////

	sh.HandleFunc("/announcements", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				announcementHandler.GetAnnouncements(w, r)
				return
			}
		case http.MethodPost:
			{
				announcementHandler.GetAddAnnouncement(w, r)
				return
			}
		case http.MethodPatch:
			{
				announcementHandler.GetUpdateAnnouncement(w, r)
				return
			}
		case http.MethodDelete:
			{
				announcementHandler.GetDeleteAnnouncement(w, r)
				return
			}
		}
	})

	sh.HandleFunc("/announcements/socket", announcementHandler.GetConnectToAnnouncementSocket)

	/////////////////////////////////////////////////////////////////////////////////

	////////////////////////////// MESSAGE ROUTER ///////////////////////////////////

	sh.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				messageHandler.GetUserMessages(w, r)
				return
			}
		case http.MethodPost:
			{
				messageHandler.GetAddMessage(w, r)
				return
			}
		}
	})

	sh.HandleFunc("/messages/read", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			{
				messageHandler.GetReadMessage(w, r)
				return
			}
		}
	})

	// sh.Handle("/messages/socket", websocket.Handler(messageHandler.GetConnectToMessageService))

	sh.HandleFunc("/messages/socket", messageHandler.GetConnectToMessageService)

	//////////////////////////////////////////////////////////////////////////////

	////////////////////////////// CERTIFICATE ROUTER ///////////////////////////////////

	sh.HandleFunc("/certificate", func(rw http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			certificateHandler.GetCertificate(rw, r)
			return
		case http.MethodGet:
			certificateHandler.GetCertificate(rw, r)
			return
		}
	})

	//1000015823856

	server := http.Server{
		Addr: ":5002",
		// Addr:         ":" + os.Getenv("PORT"),
		Handler:      accessControl(sh),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		log.Println("Server running on port", server.Addr)
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

	sig := <-sigChan
	log.Println("Terminating server : ", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(ctx)

	////////////////////////////////////////////////////////////////////////////////
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS, POST, PATCH, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}

// func setupCORS(w *http.ResponseWriter, r *http.Request) {
// 	(*w).Header().Set("Access-Control-Allow-Origin", "*")
// 	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
// 	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
// }

// func indexHandler(w http.ResponseWriter, r *http.Request) {
// 	setupCORS(&w, r)
// 	if (*r).Method == "OPTIONS" {
// 		return
// 	}
// 	// process the request...
// }
