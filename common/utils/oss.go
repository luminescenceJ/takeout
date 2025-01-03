package utils

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"takeout/global"
	"time"
)

func LocalOss(fileName string, file *multipart.FileHeader) (string, error) {
	fileData, err := file.Open()
	var (
		path string
		ext  string
	)
	defer fileData.Close()
	if err != nil {
		return "", err
	}
	path = global.Path + "/" + fileName

	// todo: change the code

	// 生成一个新的本地文件名，避免冲突
	ext = filepath.Ext(fileName)                                         // 获取文件扩展名
	newFileName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileName) // 可按需生成文件名

	// 创建目标文件
	outFile, err := os.Create(newFileName)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	// 将上传的文件内容复制到目标文件
	_, err = io.Copy(outFile, fileData)
	if err != nil {
		return "", err
	}

	// 返回文件保存的路径
	//return newFileName, nil

	return "https://www.baidu.com/img/PCtm_d9c8750bed0b3c7d089fa7d55720d6cf.png", nil
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
