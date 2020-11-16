package e2e

import (
	"os/exec"
	"strings"
	"testing"

	"github.com/MakeNowJust/heredoc/v2"
	"github.com/stretchr/testify/assert"
)

func Test_CLI(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping due to -short")
	}

	bincli, binsrv := withBinaries(t)

	t.Run("users-cli list", func(t *testing.T) {

		t.Run("should return all 30 lines of sample data", func(t *testing.T) {
			addr := "127.0.0.1:" + freePort()
			srv := startWith(t, exec.Command(binsrv, "--address", addr))
			eventuallyEqual(t, "listening", srv.Output) // Wait until listening.

			cli := startWith(t, exec.Command(bincli, "--color=never", "--address", addr, "list")).Wait()

			output := contents(cli.Output)
			assert.Equal(t, 30, strings.Count(output, "\n"))
			assert.Contains(t, output, "Wilkerson Mosley <wilkerson.mosley@email.biz> (48 years old, address: 734 Kosciusko Street, Marbury, Connecticut, 3037)")
			assert.Contains(t, output, "Alford Cole <alford.cole@email.net> (33 years old, address: 763 Halleck Street, Elbert, Nevada, 3291)")

			assert.Equal(t, 0, cli.ProcessState.ExitCode())
		})
	})

	t.Run("users-cli create", func(t *testing.T) {
		t.Run("should allow creating a user with just an email", func(t *testing.T) {
			addr := "127.0.0.1:" + freePort()
			srv := startWith(t, exec.Command(binsrv, "--address", addr))
			eventuallyEqual(t, "listening", srv.Output) // Wait until listening.

			cli := startWith(t, exec.Command(bincli, "--color=never", "--address", addr, "create", "--email=foo@bar.com")).Wait()
			assert.Equal(t, 0, cli.ProcessState.ExitCode())

			cli2 := startWith(t, exec.Command(bincli, "--color=never", "--address", addr, "list")).Wait()
			assert.Equal(t, 0, cli2.ProcessState.ExitCode())
			output := contents(cli2.Output)
			assert.Equal(t, 31, strings.Count(output, "\n"))
			assert.Contains(t, output, "  <foo@bar.com> (0 years old, address: )")
		})

		t.Run("should allow creating a user with all information", func(t *testing.T) {
			addr := "127.0.0.1:" + freePort()
			srv := startWith(t, exec.Command(binsrv, "--address", addr))
			eventuallyEqual(t, "listening", srv.Output) // Wait until listening.

			cli := startWith(t, exec.Command(bincli, "--color=never", "--address", addr,
				"create", "--email=foo@bar.com", "--firstname=Foo", "--lastname=Bar", "--age=87", "--postaladdress", "1930 Movun Point, Svalbard & Jan Mayen")).Wait()
			assert.Equal(t, 0, cli.ProcessState.ExitCode())

			cli2 := startWith(t, exec.Command(bincli, "--color=never", "--address", addr, "list")).Wait()
			assert.Equal(t, 0, cli2.ProcessState.ExitCode())
			output := contents(cli2.Output)
			assert.Equal(t, 31, strings.Count(output, "\n"))
			assert.Contains(t, output, "Foo Bar <foo@bar.com> (87 years old, address: 1930 Movun Point, Svalbard & Jan Mayen)")
		})

		t.Run("should exit with 1 when creating with an existing email", func(t *testing.T) {
			addr := "127.0.0.1:" + freePort()
			srv := startWith(t, exec.Command(binsrv, "--address", addr))
			eventuallyEqual(t, "listening", srv.Output) // Wait until listening.

			cli := startWith(t, exec.Command(bincli, "--color=never", "--address", addr, "create", "--email=wilkerson.mosley@email.biz")).Wait()
			assert.Equal(t, 1, cli.ProcessState.ExitCode())
			assert.Contains(t, contents(cli.Output), "email already exists")
		})
	})

	t.Run("users-cli get", func(t *testing.T) {

		t.Run("should print the user associated with a given email", func(t *testing.T) {
			addr := "127.0.0.1:" + freePort()
			srv := startWith(t, exec.Command(binsrv, "--address", addr))
			eventuallyEqual(t, "listening", srv.Output) // Wait until listening.

			cli := startWith(t, exec.Command(bincli, "--color=never", "--address", addr, "get", "rice.pierce@email.com")).Wait()
			assert.Equal(t, 0, cli.ProcessState.ExitCode())
			assert.Equal(t, "Rice Pierce <rice.pierce@email.com> (46 years old, address: 291 Boardwalk , Chloride, North Carolina, 8401)\n", contents(cli.Output))
		})

		t.Run("should exit with 1 when the email is not found", func(t *testing.T) {
			addr := "127.0.0.1:" + freePort()
			srv := startWith(t, exec.Command(binsrv, "--address", addr))
			eventuallyEqual(t, "listening", srv.Output) // Wait until listening.

			cli := startWith(t, exec.Command(bincli, "--color=never", "--address", addr, "get", "impossible.name@email.com")).Wait()
			assert.Equal(t, 1, cli.ProcessState.ExitCode())
			assert.Contains(t, contents(cli.Output), "email cannot be found")
		})
	})

	t.Run("users-cli search", func(t *testing.T) {
		t.Run("should print users using a part of their name", func(t *testing.T) {
			addr := "127.0.0.1:" + freePort()
			srv := startWith(t, exec.Command(binsrv, "--address", addr))
			eventuallyEqual(t, "listening", srv.Output) // Wait until listening.

			cli := startWith(t, exec.Command(bincli, "--color=never", "--address", addr, "search", "--name", "Pierce")).Wait()
			assert.Equal(t, 0, cli.ProcessState.ExitCode())
			assert.Equal(t, "Rice Pierce <rice.pierce@email.com> (46 years old, address: 291 Boardwalk , Chloride, North Carolina, 8401)\n", contents(cli.Output))
		})

		t.Run("should print users that are between two ages", func(t *testing.T) {
			addr := "127.0.0.1:" + freePort()
			srv := startWith(t, exec.Command(binsrv, "--address", addr))
			eventuallyEqual(t, "listening", srv.Output) // Wait until listening.

			cli := startWith(t, exec.Command(bincli, "--color=never", "--address", addr, "search", "--agefrom=46", "--ageto=48")).Wait()
			assert.Equal(t, 0, cli.ProcessState.ExitCode())
			assert.Equal(t, heredoc.Doc(`
				Rice Pierce <rice.pierce@email.com> (46 years old, address: 291 Boardwalk , Chloride, North Carolina, 8401)
				Angeline Stokes <angeline.stokes@email.biz> (48 years old, address: 526 Java Street, Hailesboro, Pennsylvania, 1648)
				Pacheco Fitzgerald <pacheco.fitzgerald@email.name> (48 years old, address: 278 McKibben Street, Nicholson, South Dakota, 3793)
				Wilkerson Mosley <wilkerson.mosley@email.biz> (48 years old, address: 734 Kosciusko Street, Marbury, Connecticut, 3037)
				`), contents(cli.Output))
		})

		t.Run("should print nothing and exit with 0 when no user is found", func(t *testing.T) {
			addr := "127.0.0.1:" + freePort()
			srv := startWith(t, exec.Command(binsrv, "--address", addr))
			eventuallyEqual(t, "listening", srv.Output) // Wait until listening.

			cli := startWith(t, exec.Command(bincli, "--color=never", "--address", addr, "search", "--name=Foo")).Wait()
			assert.Equal(t, 0, cli.ProcessState.ExitCode())
			assert.Contains(t, contents(cli.Output), "")
		})
	})
}
