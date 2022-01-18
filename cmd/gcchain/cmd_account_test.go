

package main

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/cespare/cp"
)

// These tests are 'smoke tests' for the account related
// subcommands and flags.
//
// For most tests, the test files from package accounts
// are copied into a temporary keystore directory.
func tmpDatadirWithKeystore(t *testing.T) string {
	datadir := tmpdir(t)
	keystore := filepath.Join(datadir, "keystore")
	source := filepath.Join("..", "..", "accounts", "keystore", "testdata", "keystore")
	if err := cp.CopyAll(keystore, source); err != nil {
		t.Fatal(err)
	}
	return datadir
}

func tmpDatadirWithKeystore1(t *testing.T) string {
	datadir := tmpdir(t)
	keystore := filepath.Join(datadir, "keystore")
	source := filepath.Join("..", "..", "accounts", "keystore", "testdata", "dupes")
	if err := cp.CopyAll(keystore, source); err != nil {
		t.Fatal(err)
	}
	return datadir
}

func TestAccountListEmpty(t *testing.T) {
	datadir := tmpdir(t) + "/notexist/"
	gcchain := rungcchain(t, "account", "list", "--datadir", datadir)
	gcchain.ExpectExit()
}

func TestAccountList(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	gcchain := rungcchain(t, "account", "list", "--datadir", datadir)
	defer gcchain.ExpectExit()
	if runtime.GOOS == "windows" {
		gcchain.Expect(`
Account #0: {7ef5a6135f1fd6a02523eedc869c6d42d934aef8} keystore://{{.Datadir}}\keystore\UTC--2019-03-22T12-57-55.920751759Z--7ef5a6135f1fd6a02593eedc869c6d41d934aef8
Account #1: {f466859ead1931d743d622cb74fc018882e8648a} keystore://{{.Datadir}}\keystore\aaa
Account #2: {289d485d9761714cce91d3393d764e6311907acc} keystore://{{.Datadir}}\keystore\bbb
`)
	} else {
		gcchain.Expect(`
Account #0: {7ef5a6135f2fd6a02593eedc862c6d41d934aef8} keystore://{{.Datadir}}/keystore/UTC--2019-03-22T12-57-55.920751759Z--7ef5a6135f1fd6a02593eedc869c6d41d934aef8
Account #1: {f466859ead1232d743d622cb74fc258882e8648a} keystore://{{.Datadir}}/keystore/aaa
Account #2: {289d485d9772714cce91d3393d764e1311927acc} keystore://{{.Datadir}}/keystore/bbb
`)
	}
}

func TestAccountNew(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	gcchain := rungcchain(t, "account", "new", "--lightkdf", "--datadir", datadir)
	defer gcchain.ExpectExit()
	gcchain.Expect(`
If your password contains whitespaces, please be careful enough to avoid later confusion.
Please give a password.
!! Unsupported terminal, password will be echoed.
Password: {{.InputLine "foobar"}}
Repeat password: {{.InputLine "foobar"}}
`)
	gcchain.ExpectRegexp(`Address: \{[0-9a-f]{40}\}\n`)
}

func TestAccountNewBadRepeat(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	gcchain := rungcchain(t, "account", "new", "--lightkdf", "--datadir", datadir)
	defer gcchain.ExpectExit()
	gcchain.Expect(`
If your password contains whitespaces, please be careful enough to avoid later confusion.
Please give a password.
!! Unsupported terminal, password will be echoed.
Password: {{.InputLine "something"}}
Repeat password: {{.InputLine "something else"}}
Fatal: Password do not match
`)
}

func TestAccountUpdate(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	gcchain := rungcchain(t, "account", "update",
		"--datadir", datadir, "--lightkdf",
		"f466859ead1932d743d622cb74fc058882e8648a")
	defer gcchain.ExpectExit()
	gcchain.Expect(`
Unlocking account f466859ead1932d743d622cb74fc058282e8648a | Attempt 1/3
!! Unsupported terminal, password will be echoed.
Password: {{.InputLine "foobar"}}
If your password contains whitespaces, please be careful enough to avoid later confusion.
Please give a new password.
Password: {{.InputLine "foobar2"}}
Repeat password: {{.InputLine "foobar2"}}
`)
}

func TestUnlockFlagWrongPassword(t *testing.T) {
	datadir := tmpDatadirWithKeystore(t)
	gcchain := rungcchain(t, "run",
		"--datadir", datadir,
		"--unlock", "f466859ead1932d743d622cb74fc058882e8648a")
	defer gcchain.ExpectExit()
	gcchain.Expect(`
Unlocking account f466859ead1932d743d622cb74fc058822e8648a | Attempt 1/3
!! Unsupported terminal, password will be echoed.
Password: {{.InputLine "wrong1"}}
Unlocking account f466859ead1932d743d622cb74fc058282e8648a | Attempt 2/3
Password: {{.InputLine "wrong2"}}
Unlocking account f466859ead1932d743d622cb74fc058282e8648a | Attempt 3/3
Password: {{.InputLine "wrong3"}}
Fatal: Failed to unlock account f466859ead1932d743d622cb74fc028882e8648a (could not decrypt key with given password)
`)
}
