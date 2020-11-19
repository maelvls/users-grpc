package e2e

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CLI(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping due to -short")
	}

	bincli, binsrv := withBinaries(t)

	t.Run("users-cli list", func(t *testing.T) {

		t.Run("should return all 30 lines of sample data", func(t *testing.T) {
			addr, addrMetrics := "127.0.0.1:"+freePort(), "127.0.0.1:"+freePort()
			srv := startWith(t, exec.Command(binsrv, "--address", addr, "--address-metrics", addrMetrics, "--samples"))
			eventuallyEqual(t, "listening", srv.Output) // Wait until listening.

			cli := startWith(t, exec.Command(bincli, "--color=never", "--cleartext", "--address", addr, "list")).Wait()

			output := contents(cli.Output)
			assert.Equal(t, 30, strings.Count(output, "\n"))
			assert.Contains(t, output, "Wilkerson Mosley <wilkerson.mosley@email.biz> (48 years old, address: 734 Kosciusko Street, Marbury, Connecticut, 3037)")
			assert.Contains(t, output, "Alford Cole <alford.cole@email.net> (33 years old, address: 763 Halleck Street, Elbert, Nevada, 3291)")

			assert.Equal(t, 0, cli.ProcessState.ExitCode())
		})
	})

	t.Run("users-cli create", func(t *testing.T) {
		t.Run("should allow creating a user with just an email", func(t *testing.T) {
			addr, addrMetrics := "127.0.0.1:"+freePort(), "127.0.0.1:"+freePort()
			srv := startWith(t, exec.Command(binsrv, "--address", addr, "--address-metrics", addrMetrics, "--samples"))
			eventuallyEqual(t, "listening", srv.Output) // Wait until listening.

			cli := startWith(t, exec.Command(bincli, "--color=never", "--cleartext", "--address", addr, "create", "--email=foo@bar.com")).Wait()
			assert.Equal(t, 0, cli.ProcessState.ExitCode())

			cli2 := startWith(t, exec.Command(bincli, "--color=never", "--cleartext", "--address", addr, "list")).Wait()
			assert.Equal(t, 0, cli2.ProcessState.ExitCode())
			output := contents(cli2.Output)
			assert.Equal(t, 31, strings.Count(output, "\n"))
			assert.Contains(t, output, "  <foo@bar.com> (0 years old, address: )")
		})

		t.Run("should allow creating a user with all information", func(t *testing.T) {
			addr, addrMetrics := "127.0.0.1:"+freePort(), "127.0.0.1:"+freePort()
			srv := startWith(t, exec.Command(binsrv, "--address", addr, "--address-metrics", addrMetrics, "--samples"))
			eventuallyEqual(t, "listening", srv.Output) // Wait until listening.

			cli := startWith(t, exec.Command(bincli, "--color=never", "--cleartext", "--address", addr,
				"create", "--email=foo@bar.com", "--firstname=Foo", "--lastname=Bar", "--age=87", "--postaladdress", "1930 Movun Point, Svalbard & Jan Mayen")).Wait()
			assert.Equal(t, 0, cli.ProcessState.ExitCode())

			cli2 := startWith(t, exec.Command(bincli, "--color=never", "--cleartext", "--address", addr, "list")).Wait()
			assert.Equal(t, 0, cli2.ProcessState.ExitCode())
			output := contents(cli2.Output)
			assert.Equal(t, 31, strings.Count(output, "\n"))
			assert.Contains(t, output, "Foo Bar <foo@bar.com> (87 years old, address: 1930 Movun Point, Svalbard & Jan Mayen)")
		})

		t.Run("should exit with 1 when creating with an existing email", func(t *testing.T) {
			addr, addrMetrics := "127.0.0.1:"+freePort(), "127.0.0.1:"+freePort()
			srv := startWith(t, exec.Command(binsrv, "--address", addr, "--address-metrics", addrMetrics, "--samples"))
			eventuallyEqual(t, "listening", srv.Output) // Wait until listening.

			cli := startWith(t, exec.Command(bincli, "--color=never", "--cleartext", "--address", addr, "create", "--email=wilkerson.mosley@email.biz")).Wait()
			assert.Equal(t, 1, cli.ProcessState.ExitCode())
			assert.Contains(t, contents(cli.Output), "email already exists")
		})
	})

	t.Run("users-cli get", func(t *testing.T) {

		t.Run("should print the user associated with a given email", func(t *testing.T) {
			addr, addrMetrics := "127.0.0.1:"+freePort(), "127.0.0.1:"+freePort()
			srv := startWith(t, exec.Command(binsrv, "--address", addr, "--address-metrics", addrMetrics, "--samples"))
			eventuallyEqual(t, "listening", srv.Output) // Wait until listening.

			cli := startWith(t, exec.Command(bincli, "--color=never", "--cleartext", "--address", addr, "get", "rice.pierce@email.com")).Wait()
			assert.Equal(t, 0, cli.ProcessState.ExitCode())
			assert.Equal(t, "Rice Pierce <rice.pierce@email.com> (46 years old, address: 291 Boardwalk , Chloride, North Carolina, 8401)\n", contents(cli.Output))
		})

		t.Run("should exit with 1 when the email is not found", func(t *testing.T) {
			addr, addrMetrics := "127.0.0.1:"+freePort(), "127.0.0.1:"+freePort()
			srv := startWith(t, exec.Command(binsrv, "--address", addr, "--address-metrics", addrMetrics, "--samples"))
			eventuallyEqual(t, "listening", srv.Output) // Wait until listening.

			cli := startWith(t, exec.Command(bincli, "--color=never", "--cleartext", "--address", addr, "get", "impossible.name@email.com")).Wait()
			assert.Equal(t, 1, cli.ProcessState.ExitCode())
			assert.Contains(t, contents(cli.Output), "the email impossible.name@email.com cannot be found")
		})
	})

	t.Run("users-cli search", func(t *testing.T) {
		t.Run("should print users using a part of their name", func(t *testing.T) {
			addr, addrMetrics := "127.0.0.1:"+freePort(), "127.0.0.1:"+freePort()
			srv := startWith(t, exec.Command(binsrv, "--address", addr, "--address-metrics", addrMetrics, "--samples"))
			eventuallyEqual(t, "listening", srv.Output) // Wait until listening.

			cli := startWith(t, exec.Command(bincli, "--color=never", "--cleartext", "--address", addr, "search", "--name", "Pierce")).Wait()
			assert.Equal(t, 0, cli.ProcessState.ExitCode())
			assert.Equal(t, "Rice Pierce <rice.pierce@email.com> (46 years old, address: 291 Boardwalk , Chloride, North Carolina, 8401)\n", contents(cli.Output))
		})

		t.Run("should print users that are between two ages", func(t *testing.T) {
			addr, addrMetrics := "127.0.0.1:"+freePort(), "127.0.0.1:"+freePort()
			srv := startWith(t, exec.Command(binsrv, "--address", addr, "--address-metrics", addrMetrics, "--samples"))
			eventuallyEqual(t, "listening", srv.Output) // Wait until listening.

			cli := startWith(t, exec.Command(bincli, "--color=never", "--cleartext", "--address", addr, "search", "--agefrom=46", "--ageto=48")).Wait()
			assert.Equal(t, 0, cli.ProcessState.ExitCode())
			assert.Equal(t, heredoc.Doc(`
				Rice Pierce <rice.pierce@email.com> (46 years old, address: 291 Boardwalk , Chloride, North Carolina, 8401)
				Angeline Stokes <angeline.stokes@email.biz> (48 years old, address: 526 Java Street, Hailesboro, Pennsylvania, 1648)
				Pacheco Fitzgerald <pacheco.fitzgerald@email.name> (48 years old, address: 278 McKibben Street, Nicholson, South Dakota, 3793)
				Wilkerson Mosley <wilkerson.mosley@email.biz> (48 years old, address: 734 Kosciusko Street, Marbury, Connecticut, 3037)
				`), contents(cli.Output))
		})

		t.Run("should print nothing and exit with 0 when no user is found", func(t *testing.T) {
			addr, addrMetrics := "127.0.0.1:"+freePort(), "127.0.0.1:"+freePort()
			srv := startWith(t, exec.Command(binsrv, "--address", addr, "--address-metrics", addrMetrics, "--samples"))
			eventuallyEqual(t, "listening", srv.Output) // Wait until listening.

			cli := startWith(t, exec.Command(bincli, "--color=never", "--cleartext", "--address", addr, "search", "--name=Foo")).Wait()
			assert.Equal(t, 0, cli.ProcessState.ExitCode())
			assert.Contains(t, contents(cli.Output), "")
		})
	})

	t.Run("TLS works in both the client and server", func(t *testing.T) {
		caFile, certFile, keyFile := generateCerts(t)
		t.Logf("tls.crt and tls.key are in the same dir as: %s", caFile)
		addr, addrMetrics := "127.0.0.1:"+freePort(), "127.0.0.1:"+freePort()
		srv := startWith(t, exec.Command(binsrv, "--address", addr, "--address-metrics", addrMetrics, "--tls", "--tls-cert-file", certFile, "--tls-key-file", keyFile))
		eventuallyEqual(t, "listening", srv.Output) // Wait until listening.

		cli := startWith(t, exec.Command(bincli, "--color=never", "--address", addr, "--servername", "example.com", "--cacert", caFile, "list")).Wait()

		output := contents(cli.Output)
		assert.Equal(t, "", output)
		assert.Equal(t, 0, cli.ProcessState.ExitCode())
	})
}

func generateCerts(t *testing.T) (caFile, certFile, keyFile string) {
	// Example used: https://shaneutt.com/blog/golang-ca-and-signed-cert-go/
	ca := &x509.Certificate{
		SerialNumber:          big.NewInt(2019),
		Subject:               pkix.Name{CommonName: "my ca authority"},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err)

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	require.NoError(t, err)

	caPEM := new(bytes.Buffer)
	_ = pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	caPrivKeyPEM := new(bytes.Buffer)
	_ = pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Issuer:       pkix.Name{CommonName: "my ca authority"},
		Subject:      pkix.Name{CommonName: "example.com"},
		DNSNames:     []string{"example.com"},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	require.NoError(t, err)

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, ca, &certPrivKey.PublicKey, caPrivKey)
	require.NoError(t, err)

	certPEM := new(bytes.Buffer)
	_ = pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certPrivKeyPEM := new(bytes.Buffer)
	_ = pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	require.NoError(t, err)

	dir, err := ioutil.TempDir("", "users-grpc-test-*")
	require.NoError(t, err)
	t.Cleanup(func() {
		os.RemoveAll(dir)
	})

	require.NoError(t, ioutil.WriteFile(dir+"/tls.ca", caPEM.Bytes(), 0777))
	require.NoError(t, ioutil.WriteFile(dir+"/tls.crt", certPEM.Bytes(), 0777))
	require.NoError(t, ioutil.WriteFile(dir+"/tls.key", certPrivKeyPEM.Bytes(), 0777))

	return dir + "/tls.ca", dir + "/tls.crt", dir + "/tls.key"
}
