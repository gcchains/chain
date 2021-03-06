

package build

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/gcchains/chain/commons/log"
)

var TestDirPrefix = "github.com/gcchains/chain/"
var DryRunFlag = flag.Bool("n", false, "dry run, don't execute commands")

// MustRun executes the given command and exits the host process for
// any error.
func MustRun(cmd *exec.Cmd) {
	fmt.Println(">>>", strings.Join(cmd.Args, " "))
	if !*DryRunFlag {
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			log.Fatal(err.Error())
		}
	}
}

func MustRunCommand(cmd string, args ...string) {
	MustRun(exec.Command(cmd, args...))
}

// GOPATH returns the value that the GOPATH environment
// variable should be set to.
func GOPATH() string {
	if os.Getenv("GOPATH") == "" {
		log.Fatal("GOPATH is not set")
	}
	return os.Getenv("GOPATH")
}

// VERSION returns the content of the VERSION file.
func VERSION() string {
	version, err := ioutil.ReadFile("VERSION")
	if err != nil {
		log.Fatal(err.Error())
	}
	return string(bytes.TrimSpace(version))
}

var warnedAboutGit bool

// RunGit runs a git subcommand and returns its output.
// The command must complete successfully.
func RunGit(args ...string) string {
	cmd := exec.Command("git", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if err := cmd.Run(); err == exec.ErrNotFound {
		if !warnedAboutGit {
			log.Println("Warning: can't find 'git' in PATH")
			warnedAboutGit = true
		}
		return ""
	} else if err != nil {
		log.Fatal(strings.Join(cmd.Args, " "), ": ", err, "\n", stderr.String())
	}
	return strings.TrimSpace(stdout.String())
}

// readGitFile returns content of file in .git directory.
func readGitFile(file string) string {
	content, err := ioutil.ReadFile(path.Join(".git", file))
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(content))
}

// Render renders the given template file into outputFile.
func Render(templateFile, outputFile string, outputPerm os.FileMode, x interface{}) {
	tpl := template.Must(template.ParseFiles(templateFile))
	render(tpl, outputFile, outputPerm, x)
}

// RenderString renders the given template string into outputFile.
func RenderString(templateContent, outputFile string, outputPerm os.FileMode, x interface{}) {
	tpl := template.Must(template.New("").Parse(templateContent))
	render(tpl, outputFile, outputPerm, x)
}

func render(tpl *template.Template, outputFile string, outputPerm os.FileMode, x interface{}) {
	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		log.Fatal(err.Error())
	}
	out, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY|os.O_EXCL, outputPerm)
	if err != nil {
		log.Fatal(err.Error())
	}
	if err := tpl.Execute(out, x); err != nil {
		log.Fatal(err.Error())
	}
	if err := out.Close(); err != nil {
		log.Fatal(err.Error())
	}
}

// CopyFile copies a file.
func CopyFile(dst, src string, mode os.FileMode) {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		log.Fatal(err.Error())
	}
	destFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer destFile.Close()

	srcFile, err := os.Open(src)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer srcFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		log.Fatal(err.Error())
	}
}

// GoTool returns the command that runs a go tool. This uses go from GOROOT instead of PATH
// so that go commands executed by build use the same version of Go as the 'host' that runs
// build code. e.g.
//
//     /usr/lib/go-1.8/bin/go run build/ci.go ...
//
// runs using go 1.8 and invokes go 1.8 tools from the same GOROOT. This is also important
// because runtime.Version checks on the host should match the tools that are run.
func GoTool(tool string, args ...string) *exec.Cmd {
	args = append([]string{tool}, args...)
	return exec.Command(filepath.Join(runtime.GOROOT(), "bin", "go"), args...)
}

// ExpandPackagesNoVendor expands a cmd/go import path pattern, skipping
// vendored packages.
func ExpandPackagesNoVendor(patterns []string) []string {
	expand := false
	for _, pkg := range patterns {
		if strings.Contains(pkg, "...") {
			expand = true
		}
	}
	if expand {
		cmd := GoTool("list", patterns...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Fatalf("package listing failed: %v\n%s", err, string(out))
		}
		var packages []string
		for _, line := range strings.Split(string(out), "\n") {
			if !strings.Contains(line, "/vendor/") {
				packages = append(packages, strings.TrimSpace(line))
			}
		}
		return packages
	}
	return patterns
}

// Read package list from build/race_test_dirs
func ReadExcludedPackagesList(excludePath string) []string {
	// TestDirPrefix
	packages := []string{}
	fi, err := os.Open(excludePath)
	defer fi.Close()

	// Exclude package set
	excludeSet := map[string]bool{}
	if err != nil {
		log.Println("Warning: can't find " + excludePath)
	} else {
		br := bufio.NewReader(fi)
		for {
			packageName, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			packageNameString := string(packageName)
			line := strings.TrimSpace(packageNameString)
			if line == "" || line[0] == '#' {
				continue
			}
			excludeSet[TestDirPrefix+packageNameString] = true
		}
	}

	// filter exclude path
	cmd := GoTool("list", "./...")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("package listing failed: %v\n%s", err, string(out))
	}
	for _, line := range strings.Split(string(out), "\n") {
		_, isExcluded := excludeSet[line]
		if isExcluded {
			continue
		}
		packages = append(packages, strings.TrimSpace(line))
	}
	return packages
}
