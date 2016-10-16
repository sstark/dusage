package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func du(path string) (stdout []byte, stderr []byte, err error) {
	cmdName := "du"
	args := []string{"-s", "-k", path}
	cmd := exec.Command(cmdName, args...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err = cmd.Run()
	stdout = outb.Bytes()
	stderr = errb.Bytes()
	return
}

func getsize(path string) (size int, err error) {
	var stdout, stderr []byte
	stdout, stderr, err = du(path)
	if err != nil {
		err = fmt.Errorf("du returned %s:\n%s", err, (string(stderr)))
		return
	}
	outs := strings.Split(string(stdout), "\t")[0]
	if outs == "" {
		err = fmt.Errorf("could not get size for %s", path)
		return
	}
	size, err = strconv.Atoi(outs)
	return
}

func getdirs(path string) (dirs []string, err error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	for _, f := range files {
		if f.Mode().IsDir() {
			dirs = append(dirs, f.Name())
		}
	}
	return
}

type dirinfo struct {
	size int
	name string
}
type dirlist []dirinfo
type dirinfoBySize dirlist

func (di dirinfo) String() string {
	return fmt.Sprintf("%s\t%s", humanBytes(float64(di.size)), di.name)
}

func (s dirinfoBySize) Len() int {
	return len(s)
}

func (s dirinfoBySize) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s dirinfoBySize) Less(i, j int) bool {
	return s[i].size < s[j].size
}

func humanBytes(b float64) (out string) {
	units := []string{"P", "T", "G", "M", "K"}
	m := len(units)
	for v, u := range units {
		fact := math.Pow(1024, float64((m - v - 1)))
		out = fmt.Sprintf("%6.2f%s", b/fact, u)
		if b > fact {
			return
		}
	}
	return
}

func main() {
	var p string
	var dl dirlist
	if len(os.Args) > 1 {
		p = os.Args[1]
	} else {
		p = "."
	}
	dirs, err := getdirs(p)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if len(dirs) == 0 {
		fmt.Fprintln(os.Stderr, "no directories found")
		os.Exit(2)
	}
	fmt.Println("getting sizes...")
	for _, dir := range dirs {
		fullpath := filepath.Join(p, dir)
		size, err := getsize(fullpath)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(3)
		} else {
			dl = append(dl, dirinfo{size, fullpath})
		}
	}
	sort.Sort(dirinfoBySize(dl))
	for _, d := range dl {
		fmt.Println(d)
	}
}
