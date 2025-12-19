package minio

import (
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"klms/internal/api/errors"
)

var Minio *minio.Client 


func MinioConnection() *minio.Client {

		endpoint := os.Getenv("MINIO_ENDPOINT")

		accesskey := os.Getenv("MINIO_ACCESSKEY")

		secretkey := os.Getenv("MINIO_SECRETKEY")


		minioClient, connerr := minio.New(endpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(accesskey,secretkey, ""),
			Secure: false,
		})

		if connerr != nil {
			log.Println(errors.ErrConnection)
		}

		return minioClient
}