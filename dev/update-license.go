
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"
)

var (
	// only files with these extensions will be considered
	extensions = []string{".go", ".js", ".qml"}

	// paths with any of these prefixes will be skipped
	skipPrefixes = []string{
		// boring stuff
		"vendor/", "tests/testdata/", "build/",
		// don't relicense vendored sources
		"crypto/bn256/",
		"crypto/ecies/",
		"crypto/secp256k1/curve.go",
		"crypto/sha3/",
		"log/",
		"common/bitutil/bitutil",
		// don't license generated files
		"contracts/dpos/contracts/admission/admission.go",
		"contracts/dpos/contracts/campaign/campaign.go",
		"contracts/dpos/contracts/pdash/pdash.go",
		"contracts/dpos/contracts/register/register.go",
		"contracts/dpos/contracts/rpt/rpt.go",
		"contracts/dpos/contracts/signerRegister/signerRegister.go",
	}

	// paths with this prefix are licensed as GPL. all other files are LGPL.
	gplPrefixes = []string{"cmd/"}

	// this regexp must match the entire license comment at the
	// beginning of each file.
	licenseCommentRE = regexp.MustCompile(`^//\s*(Copyright|This file is part of).*?\n(?://.*?\n)*\n*`)

	// this text appears at the start of AUTHORS
	authorsFileHeader = "# This is the official list of gcchain authors for copyright purposes.\n\n"
)

// this template generates the license comment.
// its input is an info structure.
var licenseT = template.Must(template.New("").Parse(`
// Copyright {{.Year}} The gcchain authors
// This file is part of {{.Whole false}}.
`[1:]))

type info struct {
	file string
	Year int64
}

func (i info) License() string {
	if i.gpl() {
		return "General Public License"
	}
	return "Lesser General Public License"
}

func (i info) ShortLicense() string {
	if i.gpl() {
		return "GPL"
	}
	return "LGPL"
}

func (i info) Whole(startOfSentence bool) string {
	if i.gpl() {
		return "gcchain"
	}
	if startOfSentence {
		return "The gcchain library"
	}
	return "the gcchain library"
}

func (i info) gpl() bool {
	for _, p := range gplPrefixes {
		if strings.HasPrefix(i.file, p) {
			return true
		}
	}
	return false
}

func main() {
	var (
		files = getFiles()
		filec = make(chan string)
		infoc = make(chan *info, 20)
		wg    sync.WaitGroup
	)

	writeAuthors(files)

	go func() {
		for _, f := range files {
			filec <- f
		}
		close(filec)
	}()
	for i := runtime.NumCPU(); i >= 0; i-- {
		// getting file info is slow and needs to be parallel.
		// it traverses git history for each file.
		wg.Add(1)
		go getInfo(filec, infoc, &wg)
	}
	go func() {
		wg.Wait()
		close(infoc)
	}()
	writeLicenses(infoc)
}

func skipFile(path string) bool {
	if strings.Contains(path, "/testdata/") {
		return true
	}
	for _, p := range skipPrefixes {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

func getFiles() []string {
	cmd := exec.Command("git", "ls-tree", "-r", "--name-only", "HEAD")
	var files []string
	err := doLines(cmd, func(line string) {
		if skipFile(line) {
			return
		}
		ext := filepath.Ext(line)
		for _, wantExt := range extensions {
			if ext == wantExt {
				goto keep
			}
		}
		return
	keep:
		files = append(files, line)
	})
	if err != nil {
		log.Fatal("error getting files:", err)
	}
	return files
}

var authorRegexp = regexp.MustCompile(`\s*[0-9]+\s*(.*)`)

func gitAuthors(files []string) []string {
	cmds := []string{"shortlog", "-s", "-n", "-e", "HEAD", "--"}
	cmds = append(cmds, files...)
	cmd := exec.Command("git", cmds...)
	var authors []string
	err := doLines(cmd, func(line string) {
		m := authorRegexp.FindStringSubmatch(line)
		if len(m) > 1 {
			authors = append(authors, m[1])
		}
	})
	if err != nil {
		log.Fatalln("error getting authors:", err)
	}
	return authors
}

func readAuthors() []string {
	content, err := ioutil.ReadFile("AUTHORS")
	if err != nil && !os.IsNotExist(err) {
		log.Fatalln("error reading AUTHORS:", err)
	}
	var authors []string
	for _, a := range bytes.Split(content, []byte("\n")) {
		if len(a) > 0 && a[0] != '#' {
			authors = append(authors, string(a))
		}
	}
	// Retranslate existing authors through .mailmap.
	// This should catch email address changes.
	authors = mailmapLookup(authors)
	return authors
}

func mailmapLookup(authors []string) []string {
	if len(authors) == 0 {
		return nil
	}
	cmds := []string{"check-mailmap", "--"}
	cmds = append(cmds, authors...)
	cmd := exec.Command("git", cmds...)
	var translated []string
	err := doLines(cmd, func(line string) {
		translated = append(translated, line)
	})
	if err != nil {
		log.Fatalln("error translating authors:", err)
	}
	return translated
}

func writeAuthors(files []string) {
	merge := make(map[string]bool)
	// Add authors that Git reports as contributorxs.
	// This is the primary source of author information.
	for _, a := range gitAuthors(files) {
		merge[a] = true
	}
	// Add existing authors from the file. This should ensure that we
	// never lose authors, even if Git stops listing them. We can also
	// add authors manually this way.
	for _, a := range readAuthors() {
		merge[a] = true
	}
	// Write sorted list of authors back to the file.
	var result []string
	for a := range merge {
		result = append(result, a)
	}
	sort.Strings(result)
	content := new(bytes.Buffer)
	content.WriteString(authorsFileHeader)
	for _, a := range result {
		content.WriteString(a)
		content.WriteString("\n")
	}
	fmt.Println("writing AUTHORS")
	if err := ioutil.WriteFile("AUTHORS", content.Bytes(), 0644); err != nil {
		log.Fatalln(err)
	}
}

func getInfo(files <-chan string, out chan<- *info, wg *sync.WaitGroup) {
	for file := range files {
		stat, err := os.Lstat(file)
		if err != nil {
			fmt.Printf("ERROR %s: %v\n", file, err)
			continue
		}
		if !stat.Mode().IsRegular() {
			continue
		}
		if isGenerated(file) {
			continue
		}
		info, err := fileInfo(file)
		if err != nil {
			fmt.Printf("ERROR %s: %v\n", file, err)
			continue
		}
		out <- info
	}
	wg.Done()
}

func isGenerated(file string) bool {
	fd, err := os.Open(file)
	if err != nil {
		return false
	}
	defer fd.Close()
	buf := make([]byte, 2048)
	n, _ := fd.Read(buf)
	buf = buf[:n]
	for _, l := range bytes.Split(buf, []byte("\n")) {
		if bytes.HasPrefix(l, []byte("// Code generated")) {
			return true
		}
	}
	return false
}

// fileInfo finds the lowest year in which the given file was committed.
func fileInfo(file string) (*info, error) {
	info := &info{file: file, Year: int64(time.Now().Year())}
	cmd := exec.Command("git", "log", "--follow", "--find-renames=80", "--find-copies=80", "--pretty=format:%ai", "--", file)
	err := doLines(cmd, func(line string) {
		y, err := strconv.ParseInt(line[:4], 10, 64)
		if err != nil {
			fmt.Printf("cannot parse year: %q", line[:4])
		}
		if y < info.Year {
			info.Year = y
		}
	})
	return info, err
}

func writeLicenses(infos <-chan *info) {
	for i := range infos {
		writeLicense(i)
	}
}

func writeLicense(info *info) {
	fi, err := os.Stat(info.file)
	if os.IsNotExist(err) {
		fmt.Println("skipping (does not exist)", info.file)
		return
	}
	if err != nil {
		log.Fatalf("error stat'ing %s: %v\n", info.file, err)
	}
	content, err := ioutil.ReadFile(info.file)
	if err != nil {
		log.Fatalf("error reading %s: %v\n", info.file, err)
	}
	// Skip file that already start with license with key word :"// Copyright"
	if strings.HasPrefix(string(content), "// Copyright") {
		fmt.Println("skipping old file:", info.file)
		return
	}

	if info.Year < 2018 {
		fmt.Println("skipping old file:", info.file)
		return
	}

	// Construct new file content.
	buf := new(bytes.Buffer)
	licenseT.Execute(buf, info)
	if m := licenseCommentRE.FindIndex(content); m != nil && m[0] == 0 {
		buf.Write(content[:m[0]])
		buf.Write(content[m[1]:])
	} else {
		buf.Write(content)
	}
	// Write it to the file.
	if bytes.Equal(content, buf.Bytes()) {
		fmt.Println("skipping (no changes)", info.file)
		return
	}
	fmt.Println("writing", info.ShortLicense(), info.file)
	if err := ioutil.WriteFile(info.file, buf.Bytes(), fi.Mode()); err != nil {
		log.Fatalf("error writing %s: %v", info.file, err)
	}
}

func doLines(cmd *exec.Cmd, f func(string)) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return err
	}
	s := bufio.NewScanner(stdout)
	for s.Scan() {
		f(s.Text())
	}
	if s.Err() != nil {
		return s.Err()
	}
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("%v (for %s)", err, strings.Join(cmd.Args, " "))
	}
	return nil
}
