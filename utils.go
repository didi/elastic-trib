package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
)

// fatal prints the error's details
// then exits the program with an exit status of 1.
func fatal(err error) {
	// make sure the error is written to the logger
	logrus.Error(err)
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

// DeDuplicate make a slice with no duplicated elements.
func DeDuplicate(input []string) []string {
	if input == nil {
		return nil
	}
	result := []string{}
	internal := map[string]struct{}{}
	for _, value := range input {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}

		if _, exist := internal[value]; !exist {
			internal[value] = struct{}{}
			result = append(result, value)
		}
	}
	return result
}

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer

	// make indent four space
	err := json.Indent(&out, []byte(in), "", "    ")
	if err != nil {
		return in
	}

	return out.String()
}

func checkURLScheme(addr, scheme string) string {
	var address string

	if strings.HasPrefix(addr, scheme) {
		address = addr
	} else {
		address = fmt.Sprintf("%s://%s", scheme, addr)
	}

	return address
}

func checkIPAddr(iparray []string) error {
	for _, ip := range iparray {
		ip = strings.TrimSpace(ip)
		if ip == "" {
			continue
		}

		i := net.ParseIP(ip)
		if i == nil {
			return fmt.Errorf("%s is invalid ip addr", ip)
		}
	}

	return nil
}

// GetCurrPath get current path string
func GetCurrPath() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	ret := path[:index]
	return ret
}

// YesOrDie get yes answer
func YesOrDie(msg string) {
	fmt.Printf("%s\n(type 'yes' to accept): ", msg)

	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')

	if !strings.EqualFold(strings.TrimSpace(text), "yes") {
		logrus.Fatalf("*** Aborting...")
	}
}
