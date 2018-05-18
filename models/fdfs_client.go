package models

import (
	"fmt"
	"github.com/weilaihui/fdfs_client"
)

//根据fastdfs文件名上传文件  ---> fileid  err
func FDFSUploadByFilename(filename string) (groupName string, fileId string, err error) {
	fdfsClient, err := fdfs_client.NewFdfsClient("./conf/client.conf")
	if err != nil {
		fmt.Println("New FdfsClient error =  ", err)
		return "", "", err
	}

	uploadResponse, err := fdfsClient.UploadByFilename(filename)
	if err != nil {
		fmt.Println("UploadByFilename error ", err)
		return "", "", err
	}

	fmt.Println(uploadResponse.GroupName)
	fmt.Println(uploadResponse.RemoteFileId)

	return uploadResponse.GroupName, uploadResponse.RemoteFileId, nil
}

//跟file的 buffer上传 ---> fileid  err
func FDFSUploadByBuffer(buffer []byte, suffix string) (gourpName string, fileId string, err error) {

	fdfsClient, err := fdfs_client.NewFdfsClient("./conf/client.conf")
	if err != nil {
		fmt.Println("New FdfsClient error =  ", err)
		return "", "", err
	}

	uploadResponse, err := fdfsClient.UploadByBuffer(buffer, suffix)
	if err != nil {
		fmt.Println("UploadByFilename error ", err)
		return "", "", err
	}

	fmt.Println(uploadResponse.GroupName)
	fmt.Println(uploadResponse.RemoteFileId)

	return uploadResponse.GroupName, uploadResponse.RemoteFileId, nil
}
