package s3

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
	"testing"
)

const (
	testId     = "testId"
	testSecret = "testSecret"
	testToken  = "testToken"
)

func MockSSOSessionHandler() (*session.Session, error) {
	creds := credentials.NewStaticCredentials(testId, testSecret, testToken)
	return &session.Session{
		Config: &aws.Config{
			Credentials: creds,
		},
	}, nil
}

func TestSSOCredentialsProvider_Retrieve(t *testing.T) {
	provider := NewSSOCredentialsProvider(MockSSOSessionHandler)

	creds, err := provider.Retrieve()
	if err != nil {
		t.Fatal(err)
	}

	if creds.AccessKeyID != testId {
		t.Fatalf("expected access key id to be %q, got %q", testId, creds.AccessKeyID)
	}

	if creds.SecretAccessKey != testSecret {
		t.Fatalf("expected secret access key to be %q, got %q", testSecret, creds.SecretAccessKey)
	}

	if creds.SessionToken != testToken {
		t.Fatalf("expected session token to be %q, got %q", testToken, creds.SessionToken)
	}

	err = os.Setenv("AWS_PROFILE", "")
	if err != nil {
		t.Fatal(err)
	}

	_, err = provider.Retrieve()
	if err == nil {
		t.Fatal("expected error: AWS_PROFILE not set")
	}
}
