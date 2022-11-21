package certificate_service

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	certificate_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/certificate/domain"
	certificate_utils "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/certificate/utils"
	grade_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/grade/domain"
	grade_repo "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/grade/repo"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

type CertificateServicePort interface {
	GenerateCertificate(token, userId string) (certificate_domain.Certificate, error)
}

type CertificateServiceAdapter struct {
	log       *log.Logger
	gradeRepo grade_repo.GradeRepoPort
}

func NewCertificateServiceAdapter(log *log.Logger, gradeRepo grade_repo.GradeRepoPort) CertificateServicePort {
	return &CertificateServiceAdapter{log, gradeRepo}
}

func (service CertificateServiceAdapter) GenerateCertificate(token, userId string) (certificate_domain.Certificate, error) {
	var certificate certificate_domain.Certificate

	user, err := certificate_utils.CheckAuth(token)

	if err != nil {
		return certificate, err
	}

	if user.Type != "REGISTRY_OFFICER" && user.Type != "SUPER_ADMIN" {
		// service.log.Println("Not authorizaed")
		return certificate, errors.New("not authorized for the operation")
	}

	gradeLabels, err := service.gradeRepo.FindGradeLabels()
	if err != nil {
		return certificate, err
	}
	allgrades, err := service.gradeRepo.FindGradesByUser(userId)
	if err != nil {
		// service.log.Println("Err in finding grades")
		return certificate, err
	}

	var grades []grade_domain.Grade = make([]grade_domain.Grade, 0)

	for a := 0; a < len(allgrades); a++ {
		if allgrades[a].Sunmitted {
			grades = append(grades, allgrades[a])
		}
	}

	if len(grades) < 1 {
		return certificate, errors.New("no grades found")
	}

	//Generate Certificate Code

	headers := []string{"Course", "Credit Hr", "Total", "Grade", "Competency"}
	contents := [][]string{}

	for k := 0; k < len(grades); k++ {
		var _totGrade float64 = 0

		_totGrade = grades[k].Assessment + grades[k].Mid + grades[k].Final

		var _competency string = "Incompetent"

		if _totGrade > 50 {
			_competency = "Competent"
		}

		var grade string = "-"

		for g := 0; g < len(gradeLabels); g++ {
			if gradeLabels[g].Min <= _totGrade && _totGrade < gradeLabels[g].Max {
				grade = gradeLabels[g].Label
			}
		}

		// service.log.Println(grades[k].Course.Name)
		// service.log.Println("grades[k].Course.Name")

		contentsContent := []string{
			grades[k].Course.Name,
			strconv.Itoa(grades[k].Course.CreditHr),
			fmt.Sprintf("%.2f", _totGrade),
			grade,
			_competency,
		}
		contents = append(contents, contentsContent)

	}

	m := pdf.NewMaroto(consts.Portrait, consts.A4)

	m.Row(30, func() {
		m.Col(12, func() {
			m.Text("Student Grade Report", props.Text{
				Size:  16,
				Style: consts.Bold,
				Align: consts.Center,
				Top:   9,
			})
		})
	})

	m.Row(6, func() {
		m.Col(1, func() {
			m.Text("Name : ", props.Text{
				Size:  10,
				Align: consts.Left,
				Top:   0,
			})
		})
		m.Col(6, func() {
			m.Text(grades[0].User.FirstName+" "+grades[0].User.LastName, props.Text{
				Size:  10,
				Style: consts.Bold,
				Align: consts.Left,
				Top:   0,
			})
		})
	})

	m.Row(6, func() {
		m.Col(1, func() {
			m.Text("Class : ", props.Text{
				Size:  10,
				Align: consts.Left,
				Top:   0,
			})
		})
		m.Col(6, func() {
			m.Text(grades[0].Section.Class.Name, props.Text{
				Size:  10,
				Style: consts.Bold,
				Align: consts.Left,
				Top:   0,
			})
		})
	})

	m.Row(6, func() {
		m.Col(1, func() {
			m.Text("Section : ", props.Text{
				Size:  10,
				Align: consts.Left,
				Top:   0,
			})
		})
		m.Col(6, func() {
			m.Text(grades[0].Section.Name, props.Text{
				Size:  10,
				Style: consts.Bold,
				Align: consts.Left,
				Top:   0,
			})
		})
	})

	// m.Row(20, func() {
	// 	m.Col(1, func() {
	// 		m.Text("Year : ", props.Text{
	// 			Size:  10,
	// 			Align: consts.Left,
	// 			Top:   0,
	// 		})
	// 	})
	// 	m.Col(3, func() {
	// 		m.Text(strconv.Itoa(int(grades[0].Section.Year)), props.Text{
	// 			Size:  10,
	// 			Style: consts.Bold,
	// 			Align: consts.Left,
	// 			Top:   0,
	// 		})
	// 	})
	// })

	m.Row(100, func() {
		m.Line(1)
		m.TableList(headers, contents, props.TableList{
			HeaderProp: props.TableListContent{
				Family:    consts.Arial,
				Style:     consts.Bold,
				Size:      11.0,
				GridSizes: []uint{3, 2, 2, 2, 3},
				// Color:     color.Color{Red: 255, Green: 0, Blue: 0},
			},
			ContentProp: props.TableListContent{
				Family:    consts.Courier,
				Style:     consts.Normal,
				Size:      10.0,
				GridSizes: []uint{3, 2, 2, 2, 3},
			},
			Align: consts.Left,
			// AlternatedBackground: &color.Color{
			// 	Red:   100,
			// 	Green: 20,
			// 	Blue:  255,
			// },
			HeaderContentSpace: 5.0,
			Line:               true,
		})
	})

	_, err = os.Open("/uploads/certificates")
	if err != nil {
		err = os.MkdirAll("uploads/certificates", 0755)
		if err != nil {
			log.Println(err)
		}
	}

	m.OutputFileAndClose("./uploads/certificates/" + strings.ReplaceAll(userId, "/", "-") + ".pdf")

	certificate.File = "/uploads/certificates/" + strings.ReplaceAll(userId, "/", "-") + ".pdf"

	return certificate, nil
}
