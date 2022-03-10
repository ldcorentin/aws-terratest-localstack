package test
// github.com/gruntwork-io/terratest v0.40.6
import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/s3"
)

const LocalstackEndpoint = "http://localhost:4566"
const REGION 			 = "eu-west-1"
var s3session *s3.S3

func init() {
    s3session = s3.New(session.Must(session.NewSession(&aws.Config{
        Region:   aws.String(REGION),
        Endpoint: aws.String(LocalstackEndpoint),
		S3ForcePathStyle: aws.Bool(true),
    })))
}

func configureTerraformOptions(t *testing.T) *terraform.Options {
	return terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// The path to where our Terraform code is located
		TerraformDir: "../terraform/",
	
		// Variables to pass to our Terraform code using -var-file options
		VarFiles: []string{"varfile.tfvars"},
	
		// Disable colors in Terraform commands so its easier to parse stdout/stderr
		NoColor: true,
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
	t.Parallel()

	terraformOptions := configureTerraformOptions(t)

	// Using "defer" runs the command at the end of the test, whether the test succeeds or fails.
	// At the end of the test, run `terraform destroy` to clean up any resources that were created
	defer terraform.Destroy(t, terraformOptions)

	// This will run `terraform init` and `terraform apply` and fail the test if there are any errors
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the values of output variables
  	bucketID := terraform.Output(t, terraformOptions, "simple_id")
  	bucketARN := terraform.Output(t, terraformOptions, "simple_arn")

  	bucketVersioningID := terraform.Output(t, terraformOptions, "versionning_id")
  	bucketVersioningARN := terraform.Output(t, terraformOptions, "versionning_arn")
	
	// Output testing
	assert.Equal(t, "aws-modules-testing-s3-simple-localstack-eu-west-1", bucketID)
	assert.Equal(t, "arn:aws:s3:::aws-modules-testing-s3-simple-localstack-eu-west-1", bucketARN)

	assert.Equal(t, "aws-modules-testing-s3-versionning-localstack-eu-west-1", bucketVersioningID)
	assert.Equal(t, "arn:aws:s3:::aws-modules-testing-s3-versionning-localstack-eu-west-1", bucketVersioningARN)

	// Versioning testing
	versioning := getVersioning(bucketVersioningID)
	assert.Equal(t, *versioning.Status, "Enabled")
}
