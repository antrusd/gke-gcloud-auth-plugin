package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/antrusd/gke-gcloud-auth-plugin/pkg/conf"
	"golang.org/x/oauth2/google"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientauthv1beta1 "k8s.io/client-go/pkg/apis/clientauthentication/v1beta1"
)

var (
	gcpScopes = []string{
		"https://www.googleapis.com/auth/cloud-platform",
		"https://www.googleapis.com/auth/userinfo.email",
	}
)

func Gcp(ctx context.Context) error {
	var cred *google.Credentials

	ec := GetExecCredential()
	//use cached exec credential
	if ec != nil {
		credString := formatJSON(ec)
		_, _ = fmt.Fprint(os.Stdout, credString)
		return nil
	}
	if conf.AppCreds != "" {
		if _, err := os.Stat(conf.AppCreds); err == nil {
			data, err := os.ReadFile(conf.AppCreds)
			if err != nil {
				return err
			}
			cred, err = google.CredentialsFromJSON(ctx, data, gcpScopes...)
			if err != nil {
				return err
			}
		}
	} else {
		//create new exec credential
		var err error
		cred, err = google.FindDefaultCredentials(ctx, gcpScopes...)
		if err != nil {
			return err
		}
	}
	if cred == nil {
		return errors.New("failed finding default credentials, cred is nil")
	}
	token, err := cred.TokenSource.Token()
	if err != nil {
		return err
	}
	if token == nil {
		return errors.New("failed retrieving token from credentials")
	}
	ec = newExecCredential(token.AccessToken, token.Expiry)
	//cache exec credential
	SaveExecCredential(ec)
	credString := formatJSON(ec)
	_, _ = fmt.Fprint(os.Stdout, credString)
	return nil
}

func formatJSON(ec *clientauthv1beta1.ExecCredential) string {
	//pretty print
	enc, _ := json.MarshalIndent(ec, "", "  ")
	return string(enc)
}

func newExecCredential(token string, exp time.Time) *clientauthv1beta1.ExecCredential {
	metaExp := metav1.NewTime(exp)
	//the google token sometimes contains trailing periods,
	//they cause problems with various tools, thus right trim
	token = strings.TrimRightFunc(token, func(r rune) bool {
		return r == '.'
	})
	_ = clientauthv1beta1.SchemeGroupVersion.Identifier()
	ec := &clientauthv1beta1.ExecCredential{
		TypeMeta: metav1.TypeMeta{
			APIVersion: clientauthv1beta1.SchemeGroupVersion.Identifier(),
			Kind:       "ExecCredential",
		},
		Status: &clientauthv1beta1.ExecCredentialStatus{
			ExpirationTimestamp: &metaExp,
			Token:               token,
		},
	}
	return ec
}
