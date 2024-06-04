package tomorrowio

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
)

func LoadEnv(paths ...string) {
	if len(paths) == 0 {
		// No paths were supplied. Search the cwd.
		cwd, err := os.Getwd()
		if err != nil {
			msg := "WARNING: No paths were provided and cwd could not be determined."
			log.Println(msg)
		}
		paths = append(paths, cwd)
	}

	for _, p := range paths {
		envvars := make(map[string]string)
		dir, err := os.ReadDir(p)
		if err != nil {
			msg := fmt.Sprintf("WARNING: Could not read env path %s", p)
			log.Printf(msg)
		}

		for _, f := range dir {
			if strings.Contains(f.Name(), ".env") {
				envvars, err = readEnvFile(f.Name())
				if err != nil {
					msg := fmt.Sprintf("WARNING: Could not read env file %s", f)
					log.Printf(msg)
				}
				// Set the environment variables.
				for k, v := range envvars {
					os.Setenv(k, v)
				}
			}
		}
	}
}

func NewHTTPClient() *http.Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Got error while creating cookie jar %s", err.Error())
	}
	return &http.Client{
		Jar:           jar,
		CheckRedirect: minRedirects,
		Timeout:       DefaultHTTPTimeout,
	}
}

func minRedirects(req *http.Request, via []*http.Request) error {
	// Print the URL of the redirected request
	fmt.Println("Redirected to:", req.URL.String())
	// Allow a maximum of 5 redirects
	if len(via) >= 5 {
		return errors.New("too many redirects")
	}
	return nil
}

func readEnvFile(filename string) (map[string]string, error) {
	vars := make(map[string]string)

	filehandle, err := os.Open(filename)
	if err != nil {
		return vars, err
	}

	fileScanner := bufio.NewScanner(filehandle)
	for fileScanner.Scan() {
		lineChunks := strings.Split(fileScanner.Text(), "=")
		vars[lineChunks[0]] = lineChunks[1]
	}
	return vars, nil
}
