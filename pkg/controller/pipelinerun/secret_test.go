package pipelinerun

import (
	"regexp"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
)

func TestGetAuthSecretWithExistingToken(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	testToken := "abcdefghijklmnopqrstuvwxyz12345678901234"

	secret := &corev1.Secret{
		Type: corev1.SecretTypeOpaque,
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: testNamespace,
		},
		Data: map[string][]byte{
			"token": []byte(testToken),
		},
	}
	objs := []runtime.Object{
		secret,
	}

	cl := fake.NewFakeClient(objs...)
	sec, err := getAuthSecret(cl, getNamespaceSecretName(testNamespace))
	if err != nil {
		t.Fatal(err)
	}
	if sec != testToken {
		t.Fatalf("got %s, wanted %s", sec, testToken)
	}
}

func TestGetAuthSecretWithNoSecret(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	objs := []runtime.Object{}

	cl := fake.NewFakeClient(objs...)
	_, err := getAuthSecret(cl, getNamespaceSecretName(testNamespace))

	wantedErr := "error getting secret test-namespace/github-auth.* not found"
	if !matchError(t, wantedErr, err) {
		t.Fatalf("failed to match error when no secret: got %s, wanted %s", err, wantedErr)
	}
}

func TestGetAuthSecretWithNoToken(t *testing.T) {
	logf.SetLogger(logf.ZapLogger(true))
	secret := &corev1.Secret{
		Type: corev1.SecretTypeOpaque,
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: testNamespace,
		},
		Data: map[string][]byte{},
	}
	objs := []runtime.Object{
		secret,
	}

	cl := fake.NewFakeClient(objs...)
	_, err := getAuthSecret(cl, getNamespaceSecretName(testNamespace))

	wantedErr := "secret .* does not have a 'token' key"
	if !matchError(t, wantedErr, err) {
		t.Fatalf("failed to match error when no secret: got %s, wanted %s", err, wantedErr)
	}
}

func matchError(t *testing.T, s string, e error) bool {
	t.Helper()
	match, err := regexp.MatchString(s, e.Error())
	if err != nil {
		t.Fatal(err)
	}
	return match
}