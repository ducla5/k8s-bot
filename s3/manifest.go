package s3

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io/ioutil"
	"k8s-bot/config"
)


type s3Util struct {
	s3Client *s3.S3
}

type Util interface {
	GetManifestFromS3(key string) (string, error)
}

func NewS3Client() Util {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-northeast-1")},
	)

	if err != nil {
		panic(err)
	}

	s3Client := s3.New(sess)
	return &s3Util{s3Client: s3Client}
}

func (u *s3Util)GetManifestFromS3(key string) (string, error) {
	requestInput := &s3.GetObjectInput{
		Bucket: aws.String(config.C.S3.BukkenName),
		Key:    aws.String(key),
	}

	result, err := u.s3Client.GetObject(requestInput)

	if err != nil {
		return "", err
	}

	defer result.Body.Close()
	body1, err := ioutil.ReadAll(result.Body)

	if err != nil {
		return "", err
	}

	bodyString1 := fmt.Sprintf("%s", body1)

	return bodyString1, nil
}