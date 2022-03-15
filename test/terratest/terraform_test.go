package test

// github.com/gruntwork-io/terratest v0.40.6
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

const LocalstackEndpoint = "http://localhost:4566"
const REGION = "eu-west-1"

var s3session *s3.S3

func init() {
	s3session = s3.New(session.Must(session.NewSession(&aws.Config{
		Region:           aws.String(REGION),
		Endpoint:         aws.String(LocalstackEndpoint),
		S3ForcePathStyle: aws.Bool(true),
	})))
}

func configureTerraformOptions(t *testing.T, target string) *terraform.Options {
	return terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../terraform/",

		// Variables to pass to our Terraform code using -var-file options
		VarFiles: []string{"varfile.tfvars"},

		// Disable colors in Terraform commands so its easier to parse stdout/stderr
		NoColor: true,

		// Test targets
		Targets: []string{target},
	})
}

func getVersioning(bucketVersioningID string) (resp *s3.GetBucketVersioningOutput) {
	resp, err := s3session.GetBucketVersioning(&s3.GetBucketVersioningInput{
		Bucket: aws.String(bucketVersioningID),
	})

	if err != nil {
		panic(err)
	}

	return resp
}

func TestTerraformS3(t *testing.T) {

	terraformOptions := configureTerraformOptions(t, "module.test")

	// Using "defer" runs the command at the end of the test, whether the test succeeds or fails.
	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the values of output variables
	bucketID := terraform.Output(t, terraformOptions, "simple_id")
	bucketARN := terraform.Output(t, terraformOptions, "simple_arn")

	// Output testing
	assert.Equal(t, "aws-modules-testing-s3-simple-localstack-eu-west-1", bucketID)
	assert.Equal(t, "arn:aws:s3:::aws-modules-testing-s3-simple-localstack-eu-west-1", bucketARN)

	// Versioning testing
	nullVersioning := getVersioning(bucketID)
	assert.Equal(t, reflect.Indirect(reflect.ValueOf(*nullVersioning)).Type().Field(0).Name, "_")
}

func TestTerraformS3Versioning(t *testing.T) {

	terraformOptions := configureTerraformOptions(t, "module.test_versioning")

	// Using "defer" runs the command at the end of the test, whether the test succeeds or fails.
	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the values of output variables
	bucketVersioningID := terraform.Output(t, terraformOptions, "versioning_id")
	bucketVersioningARN := terraform.Output(t, terraformOptions, "versioning_arn")

	// Output testing
	assert.Equal(t, "aws-modules-testing-s3-versioning-localstack-eu-west-1", bucketVersioningID)
	assert.Equal(t, "arn:aws:s3:::aws-modules-testing-s3-versioning-localstack-eu-west-1", bucketVersioningARN)

	// Versioning testing
	versioning := getVersioning(bucketVersioningID)
	assert.Equal(t, *versioning.Status, "Enabled")
}

func TestTerraformS3Interactions(t *testing.T) {

	terraformOptions := configureTerraformOptions(t, "module.test")

	// Using "defer" runs the command at the end of the test, whether the test succeeds or fails.
	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the values of output variables
	bucketID := terraform.Output(t, terraformOptions, "simple_id")
	bucketARN := terraform.Output(t, terraformOptions, "simple_arn")

	// Output testing
	assert.Equal(t, "aws-modules-testing-s3-simple-localstack-eu-west-1", bucketID)
	assert.Equal(t, "arn:aws:s3:::aws-modules-testing-s3-simple-localstack-eu-west-1", bucketARN)

	// Put, Get test
	body := map[string]string{"foo": "bar"}
	file, err := json.Marshal(body)
	_, err = s3session.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketID),
		Key:    aws.String("test"),
		Body:   aws.ReadSeekCloser(bytes.NewReader(file)),
	})

	if err != nil {
		fmt.Println("Got error uploading file:")
		fmt.Println(err)
	}

	resp, err := s3session.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketID),
		Key:    aws.String("test"),
	})
	if err != nil {
		fmt.Println("Got error downloading file:")
		fmt.Println(err)
	}

	size := int(*resp.ContentLength)
	buffer := make([]byte, size)
	defer resp.Body.Close()
	var bbuffer bytes.Buffer

	for true {
		num, rerr := resp.Body.Read(buffer)
		if num > 0 {
			bbuffer.Write(buffer[:num])
		} else if rerr == io.EOF || rerr != nil {
			break
		}
	}

	respBody := map[string]string{}
	json.Unmarshal([]byte(bbuffer.String()), &respBody)
	assert.Equal(t, respBody, body)

	_, err = s3session.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucketID),
		Key:    aws.String("test"),
	})
	if err != nil {
		fmt.Println(err)
	}

	err = s3session.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucketID),
		Key:    aws.String("test"),
	})
	if err != nil {
		fmt.Println(err)
	}
}
