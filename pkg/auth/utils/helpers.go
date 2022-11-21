package auth_utils

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	auth_domain "github.com/dawitaschalew-mock/hahu_sms_backend/pkg/auth/domain"
)

const (
	BUCKET_NAME = "hahusmsbucket"
	REGION      = "us-east-1"
	PART_SIZE   = 5000000
	RETRIES     = 2
)

var (
	s3session *s3.S3
)

func initS3() {
	s3session = s3.New(session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(REGION),
		Credentials: credentials.NewStaticCredentials("AKIATLRHCAJ7OTROBI5Q", "lGHhY0IXUHtoqxYy1BjQWEHH+EXxOGNqKBDq728D", ""),
	})))
}

func StoreProfilePicFromFormFileToAWSS3(r *http.Request, user auth_domain.User) string {
	initS3()
	var storedFileAdd string = ""

	file_, fileHandler, err := r.FormFile("profilePic")
	if err == nil {
		if file_ != nil {
			name := user.Username + time.Now().UTC().String() + "." + strings.Split(fileHandler.Filename, ".")[len(strings.Split(fileHandler.Filename, "."))-1]
			fileSize := fileHandler.Size

			buffer := make([]byte, fileSize)
			file_.Read(buffer)

			expiryDate := time.Now().AddDate(0, 0, 1)

			createdRes, err := s3session.CreateMultipartUpload(&s3.CreateMultipartUploadInput{
				Bucket:  aws.String(BUCKET_NAME),
				Key:     aws.String(name),
				Expires: &expiryDate,
			})

			if err != nil {
				log.Println("err in aws")
				log.Println(err)
				return storedFileAdd
			}

			var start, currentSize int
			var remaining int = int(fileSize)
			var partNum = 1
			var completedParts []*s3.CompletedPart

			for start = 0; remaining != 0; start += PART_SIZE {

				if remaining < PART_SIZE {
					currentSize = remaining
				} else {
					currentSize = PART_SIZE
				}

				completed, err := Upload(createdRes, buffer[start:start+currentSize], partNum)
				if err != nil {
					_, err := s3session.AbortMultipartUpload(&s3.AbortMultipartUploadInput{
						Bucket:   createdRes.Bucket,
						Key:      createdRes.Key,
						UploadId: createdRes.UploadId,
					})
					if err != nil {
						log.Println("err in aws 1")
						log.Println(err)
						return storedFileAdd
					}
					log.Println("err aws 2")
					log.Println(err)
					return storedFileAdd
				}

				remaining -= currentSize
				log.Println("Uploading")

				partNum++

				completedParts = append(completedParts, completed)

			}

			res, err := s3session.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
				Bucket:   createdRes.Bucket,
				Key:      createdRes.Key,
				UploadId: createdRes.UploadId,
				MultipartUpload: &s3.CompletedMultipartUpload{
					Parts: completedParts,
				},
			})

			if err != nil {
				log.Println("err final upload aws")
				log.Println(err)
				return storedFileAdd
			}

			storedFileAdd = *res.Location
		}

		defer file_.Close()
	}

	return storedFileAdd
}

func Upload(res *s3.CreateMultipartUploadOutput, fileBytes []byte, partNum int) (completedPart *s3.CompletedPart, err error) {
	var try int
	for try <= RETRIES {
		uploadRes, err := s3session.UploadPart(&s3.UploadPartInput{
			Body:          bytes.NewReader(fileBytes),
			Bucket:        res.Bucket,
			Key:           res.Key,
			PartNumber:    aws.Int64(int64(partNum)),
			UploadId:      res.UploadId,
			ContentLength: aws.Int64(int64(len(fileBytes))),
		})

		if err != nil {
			log.Println("err aws upload")
			log.Println(err)

			if try == RETRIES {
				return nil, err
			} else {
				try++
			}
		} else {
			return &s3.CompletedPart{
				ETag:       uploadRes.ETag,
				PartNumber: aws.Int64(int64(partNum)),
			}, nil
		}

	}

	return nil, nil
}

func StoreFileFromMultiPartRequest(file multipart.File, fileName string, path string) (*os.File, error) {
	_, err := os.Open(path)
	if err != nil {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return nil, err
		}
	}

	storedFile, err := os.OpenFile("./"+path+"/"+fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	io.Copy(storedFile, file)

	return storedFile, nil
}

/*

"user": {
        "firstName": "test",
        "lastName": "test",
        "email": "test.testadmin@gmaill.com",
        "phone": "0911164081",
        "username": "testAdmin",
        "password": "testpassword",
        "profilePic": "",
        "address": {
            "country": "Ethiopia",
            "region": "Addis Ababa",
            "city": "Addis Ababa",
            "subCity": "yeka",
            "woreda": 11,
            "houseNo": 1234
        }
    }

*/
