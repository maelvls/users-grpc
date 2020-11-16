package e2e

import (
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"syscall"
	"testing"
	"time"

	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
	"github.com/phayes/freeport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Returns the path to the built CLI. Better call it only once since it
// needs to recompile.
func withBinaries(t *testing.T) (string, string) {
	start := time.Now()

	bincli, err := gexec.Build("github.com/maelvls/users-grpc/users-cli")
	require.NoError(t, err)

	binsrv, err := gexec.Build("github.com/maelvls/users-grpc/users-server")
	require.NoError(t, err)

	t.Cleanup(func() {
		gexec.Terminate()
		gexec.CleanupBuildArtifacts()
	})

	t.Logf("compiling binaries took %v, path: %s", time.Since(start).Truncate(time.Second), binsrv)
	return bincli, binsrv
}

type e2ecmd struct {
	*exec.Cmd
	Output *gbytes.Buffer // Both stdout and stderr.
	T      *testing.T
}

func (cmd *e2ecmd) Wait() *e2ecmd {
	_ = cmd.Cmd.Wait()
	return cmd
}

// Runs the passed command and make sure SIGTERM is called on cleanup. Also
// dumps stderr and stdout using log.Printf.
func startWith(t *testing.T, cmd *exec.Cmd) *e2ecmd {
	buff := gbytes.NewBuffer()
	cmd.Stdout = createWriterLoggerStr("stdout", buff)
	cmd.Stderr = createWriterLoggerStr("stderr", buff)

	err := cmd.Start()
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = cmd.Process.Signal(syscall.SIGTERM)
	})

	return &e2ecmd{Cmd: cmd, Output: buff, T: t}
}

func contents(f io.Reader) string {
	bytes, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

// createWriterLoggerStr returns a writer that behaves like w except that it
// logs (using log.Printf) each write to standard error, printing the
// prefix and the data written as a string.
//
// Pretty much the same as iotest.NewWriterLogger except it logs strings,
// not hexadecimal jibberish.
func createWriterLoggerStr(prefix string, w io.Writer) io.Writer {
	return &writeLogger{prefix, w}
}

type writeLogger struct {
	prefix string
	w      io.Writer
}

func (l *writeLogger) Write(p []byte) (n int, err error) {
	n, err = l.w.Write(p)
	if err != nil {
		log.Printf("%s %s: %v", l.prefix, string(p[0:n]), err)
	} else {
		log.Printf("%s %s", l.prefix, string(p[0:n]))
	}
	return
}

func freePort() string {
	port, err := freeport.GetFreePort()
	if err != nil {
		panic(err)
	}
	return strconv.Itoa(port)
}

// When given a io.Reader, checks that the given string eventuall appears.
// A bit like Testify's require.Eventually but works directly on a
// io.Reader.
//
// Commented out since I'm not using it anymore.

func eventuallyEqual(t *testing.T, expected string, got *gbytes.Buffer, msgsAndArgs ...interface{}) {
	expectedBuffer := gbytes.Say(expected)

	match := func() func() bool {
		return func() bool {
			ok, err := expectedBuffer.Match(got)
			assert.NoError(t, err)

			return ok
		}
	}

	if !assert.Eventually(t, match(), 2*time.Second, 100*time.Millisecond, msgsAndArgs...) {
		t.Errorf(expectedBuffer.FailureMessage(expected))
	}
}
