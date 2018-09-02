package minio

import (
	"fmt"
	"html/template"
	"log"

	//	"os"
	"time"

	"github.com/minio/minio-go"
)

type Info struct {
	Endpoint  string
	Port      int
	AccessKey string
	SecretKey string
	Secure 	  bool
}

type ObjectHandlers struct {
	minioClient *minio.Client
}

func (c Info) Connect() *ObjectHandlers {

	minioClient, err := minio.New(fmt.Sprintf("%v:%v", c.Endpoint, c.Port), c.AccessKey, c.SecretKey, c.Secure)
	if err != nil {
		log.Fatalln(err)
	}

	objectHandler := &ObjectHandlers{
		minioClient: minioClient,
	}

	return objectHandler
}

func (o *ObjectHandlers) Map(baseURI string) template.FuncMap {
	f := make(template.FuncMap)

	f["CDNCSS"] = func(bucket string, filepath string, media string) template.HTML {

		path, err := o.GetPresignedURLHandler(bucket, filepath)

		if err != nil {
			log.Println("CSS Error:", err)
			return template.HTML("<!-- CSS Error: " + path + " -->")
		}

		return template.HTML(fmt.Sprintf(`<link media="%v" rel="stylesheet" type="text/css" href="%v" />`, media, path))
	}

	f["CDNJS"] = func(bucket string, filepath string) template.HTML {

		path, err := o.GetPresignedURLHandler(bucket, filepath)

		if err != nil {
			log.Println("JS Error:", err)
			return template.HTML("<!-- JS Error: " + path + " -->")
		}

		return template.HTML(fmt.Sprintf(`<script type="text/javascript" src="%v"></script>`, path))
	}

	f["CDNIMG"] = func(bucket string, filepath string, class string) template.HTML {

		path, err := o.GetPresignedURLHandler(bucket, filepath)

		if err != nil {
			log.Println("CSS Error:", err)
			return template.HTML(fmt.Sprintf(`<img class="%v" src="img/placeholder.jpg" />`, class))
		}

		return template.HTML(fmt.Sprintf(`<img class="%v" src="%v" />`, class, path))
	}

	return f
}

// GetPresignedURLHandler - generates presigned access URL for an object.
func (o *ObjectHandlers) GetPresignedURLHandler(bucketName string, objectName string) (string, error) {
	// The object for which the presigned URL has to be generated is sent as a query
	// parameter from the client.
	var err error
	if objectName == "" {
		return "", err
	}
	
	presignedURL, err := o.minioClient.PresignedGetObject(bucketName, objectName, 30*time.Second, nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", presignedURL), err
}
