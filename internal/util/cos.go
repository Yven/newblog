package util

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/tencentyun/cos-go-sdk-v5"
)

type Cos struct {
	Client    *cos.Client
	BucketURL string
	SecretID  string
	SecretKey string
}

func NewCos(bucketURL, secretID, secretKey string) *Cos {
	return &Cos{
		BucketURL: bucketURL,
		SecretID:  secretID,
		SecretKey: secretKey,
	}
}

func (c *Cos) UploadStream(name string, file io.Reader) (string, error) {
	u, _ := url.Parse(c.BucketURL)
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  c.SecretID,
			SecretKey: c.SecretKey,
		},
	})

	name = strings.TrimPrefix(name, "/")

	// 已经存在
	ok, err := client.Object.IsExist(context.Background(), name)
	if err == nil && ok {
		return u.String() + "/" + name, nil
	}

	_, err = client.Object.Put(context.Background(), name, file, nil)
	if err != nil {
		return "", err
	}

	return u.String() + "/" + name, nil
}
