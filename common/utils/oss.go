package utils

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"mime/multipart"
	"os"
	"takeout/global"
	"time"
)

func LocalOss(fileName string, file *multipart.FileHeader) (string, error) {
	fileData, err := file.Open()
	defer fileData.Close()
	if err != nil {
		return "", err
	}

	// 时间戳 + uuid + 文件后缀
	path := global.Config.Path.LocalPath
	fileName = fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileName)
	savePath := path + "/" + fileName

	// 创建目标文件
	outFile, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	// 将上传的文件内容复制到目标文件
	if _, err = io.Copy(outFile, fileData); err != nil {
		return "", err
	}

	visitPath := global.Config.Path.VisitPath + "/" + fileName
	return visitPath, nil
}

func AliyunOss(fileName string, file *multipart.FileHeader) (string, error) {
	config := global.Config.AliOss
	client, err := oss.New(config.EndPoint, config.AccessKeyId, config.AccessKeySecret)
	if err != nil {
		return "", err
	}
	bucket, err := client.Bucket(config.BucketName)
	if err != nil {
		return "", err
	}

	fileData, err := file.Open()
	defer fileData.Close()

	err = bucket.PutObject(fileName, fileData)
	if err != nil {
		return "", err
	}
	imagePath := "https://" + config.BucketName + "." + config.EndPoint + "/" + fileName
	fmt.Println("文件上传到：", imagePath)
	return imagePath, nil
}
