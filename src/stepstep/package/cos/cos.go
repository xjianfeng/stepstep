package cos

import (
	"bytes"
	"context"
	tcos "github.com/nelsonken/cos-go-sdk-v5/cos"
	"io"
	"stepstep/conf"
	"time"
)

var (
	cosOption *tcos.Option
)

func SetCosOption() {
	cosOption = &tcos.Option{
		AppID:     conf.CfgCos.AppId,
		SecretID:  conf.CfgCos.SecretId,
		SecretKey: conf.CfgCos.SecretKey,
		Region:    conf.CfgCos.Region,
		Bucket:    conf.CfgCos.Bucket,
	}
}

func UploadObject(fileName string, content io.Reader) error {
	client := tcos.New(cosOption)
	bucket := client.Bucket(conf.CfgCos.Bucket)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := bucket.UploadObject(ctx, fileName, content, &tcos.AccessControl{})
	if err != nil {
		return err
	}
	return nil
}

func UploadBytes(fileName string, content []byte) error {
	client := tcos.New(cosOption)
	bucket := client.Bucket(conf.CfgCos.Bucket)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	b := bytes.NewReader(content)
	err := bucket.UploadObject(ctx, fileName, b, &tcos.AccessControl{})
	if err != nil {
		return err
	}
	return nil
}

func DeleteObject(fileName string) error {
	client := tcos.New(cosOption)
	bucket := client.Bucket(conf.CfgCos.Bucket)
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := bucket.DeleteObject(ctx, fileName)
	if err != nil {
		return err
	}
	return nil
}
