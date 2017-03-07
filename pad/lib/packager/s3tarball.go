package packager

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ahamilton55/fs-test/pad/utils"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

const (
	latestFilename  = "latest-build"
	BuildDateFormat = "0601021504"
)

type S3Tarball struct {
	Bucket string
	Config utils.CommandConfig
}

func (st S3Tarball) Build() (string, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	tmpDir, err := ioutil.TempDir(os.TempDir(), st.Config.Service)

	tarFilename := filepath.Join(tmpDir, fmt.Sprintf("%s.tar", st.Config.Service))
	tarfile, err := os.Create(tarFilename)
	if err != nil {
		return "", err
	}
	defer tarfile.Close()

	tarball := tar.NewWriter(tarfile)
	defer tarball.Close()

	err = filepath.Walk(pwd,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if strings.HasPrefix(strings.TrimPrefix(path, pwd), "/.") {
				return nil
			}

			header, err := tar.FileInfoHeader(info, info.Name())
			if err != nil {
				return err
			}

			header.Name = filepath.Join(".", strings.TrimPrefix(path, pwd))

			if err := tarball.WriteHeader(header); err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(tarball, file)
			return err
		})

	if err != nil {
		return "", err
	}
	err = tarball.Close()
	if err != nil {
		return "", err
	}

	reader, err := os.Open(tarFilename)
	if err != nil {
		return "", err
	}

	filename := filepath.Base(tarFilename)
	target := filepath.Join(tmpDir, fmt.Sprintf("%s-%s.tgz", time.Now().Format(BuildDateFormat), st.Config.Service))
	writer, err := os.Create(target)
	if err != nil {
		return "", err
	}
	defer writer.Close()

	archiver := gzip.NewWriter(writer)
	archiver.Name = filename
	defer archiver.Close()

	_, err = io.Copy(archiver, reader)
	return target, err
}

func (st S3Tarball) Push(filename string) error {
	sess, err := utils.GetAWSSession("", st.Config.Profile)
	if err != nil {
		return err
	}

	svc := s3.New(sess)

	key := fmt.Sprintf("%s/%s", st.Config.Service, filepath.Base(filename))

	reader, err := os.Open(filename)
	if err != nil {
		return err
	}

	input := &s3.PutObjectInput{
		Bucket: aws.String(st.Bucket),
		Key:    aws.String(key),
		Body:   reader,
	}

	_, err = svc.PutObject(input)
	if err != nil {
		return err
	}

	key = fmt.Sprintf("%s/%s", st.Config.Service, latestFilename)

	input = &s3.PutObjectInput{
		Bucket: aws.String(st.Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader([]byte(filepath.Base(filename))),
	}

	_, err = svc.PutObject(input)
	if err != nil {
		return err
	}

	return nil
}

func (st S3Tarball) Cleanup(file string) error {
	if file == "" {
		return fmt.Errorf("No file provided for cleanup")
	}

	dir := filepath.Dir(file)
	if dir == "." {
		return fmt.Errorf("Won't clean up relative directories")
	}

	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}

	return nil
}

func (st S3Tarball) FindPackage(ver string) (string, error) {
	sess, err := utils.GetAWSSession("", st.Config.Profile)
	if err != nil {
		return "", err
	}

	svc := s3.New(sess)
	getInput := &s3.GetObjectInput{
		Bucket: aws.String(st.Bucket),
	}

	if ver == "latest" {
		getInput.Key = aws.String(fmt.Sprintf("%s/%s", st.Config.Service, latestFilename))
	} else {
		getInput.Key = aws.String(fmt.Sprintf("%s/%s-%s", st.Config.Service, ver, st.Config.Service))
	}

	resp, err := svc.GetObject(getInput)
	if err != nil {
		return "", err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	objectName := strings.TrimSpace(string(content))

	listInput := &s3.ListObjectsInput{
		Bucket: aws.String(st.Bucket),
		Prefix: aws.String(fmt.Sprintf("%s/%s", st.Config.Service, objectName)),
	}

	listResp, err := svc.ListObjects(listInput)
	if err != nil {
		return "", err
	}

	if len(listResp.Contents) != 1 {
		return "", fmt.Errorf("Could not find latest build")
	}

	return fmt.Sprintf("s3://%s/%s/%s", st.Bucket, st.Config.Service, objectName), nil
}
