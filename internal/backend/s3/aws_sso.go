package s3

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"os"
)

type SSOCredentialsProvider struct {
	sessionHandler func() (*session.Session, error)
	retrieved      bool
}

func NewSSOCredentialsProvider(sessionHandler func() (*session.Session, error)) *SSOCredentialsProvider {
	return &SSOCredentialsProvider{
		sessionHandler: sessionHandler,
	}
}

func SSOSessionHandler() (*session.Session, error) {
	profile := os.Getenv("AWS_PROFILE")
	if profile == "" {
		return nil, fmt.Errorf("AWS_PROFILE not set")
	}

	return session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Profile:           profile,
	})
}

func (e *SSOCredentialsProvider) Retrieve() (credentials.Value, error) {
	e.retrieved = false

	profile := os.Getenv("AWS_PROFILE")
	if profile == "" {
		return credentials.Value{}, fmt.Errorf("AWS_PROFILE not set")
	}

	sess := session.Must(e.sessionHandler())

	creds, err := sess.Config.Credentials.Get()
	if err != nil {
		return credentials.Value{}, fmt.Errorf("failed to get credentials: %w", err)
	}

	e.retrieved = true
	return credentials.Value{
		AccessKeyID:     creds.AccessKeyID,
		SecretAccessKey: creds.SecretAccessKey,
		SessionToken:    creds.SessionToken,
	}, nil
}

func (e *SSOCredentialsProvider) IsExpired() bool {
	return !e.retrieved
}
