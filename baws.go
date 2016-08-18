package main

// note: until the go libraries can read region from ~/.aws/credentials,
// it should be set as an environment variable, i.e. export AWS_REGION=us-west-2

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	s3Command   = kingpin.Command("s3", "S3 Commands")
	s3LsCommand = s3Command.Command("ls", "List files in a bucket")
	s3LsPath    = s3LsCommand.Arg("path", "Path to list files in (s3://mybucket/foo)").String()
	include     = kingpin.Flag("include", "Include file mask").String()
)

func main() {
	switch kingpin.Parse() {
	case "s3 ls":
		s3Ls()
	}
}

func s3Ls() {
	sess, err := session.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	u, err := url.Parse(*s3LsPath)
	if err != nil {
		fmt.Println("Failed to understand path", err)
		return
	}

	u.Path = strings.TrimLeft(u.Path, "/")

	params := &s3.ListObjectsInput{
		Bucket: aws.String(u.Host),
		Prefix: aws.String(u.Path),
	}

	svc := s3.New(sess)
	err = svc.ListObjectsPages(params, printObjects)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

func printObjects(resp *s3.ListObjectsOutput, lastPage bool) bool {
	if resp.Contents == nil {
		fmt.Println("No results found.")
		return false
	}

	for _, line := range resp.Contents {
		match := true

		if *include != "" {
			match, _ = filepath.Match(*include, *line.Key)
		}

		if match {
			fmt.Println(*line.Key)
		}
	}

	return true

}
