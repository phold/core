package main

import (
	"io"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/oauth2"
)

type AccessToken string

func (t *AccessToken) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: string(*t),
	}
	return token, nil
}

// var token = flag.String("token", "", "Digital Ocean API Token")

func main() {
	/*
		flag.Parse()
		tokenSource := AccessToken(*token)
		oauthClient := oauth2.NewClient(oauth2.NoContext, &tokenSource)
		client := godo.NewClient(oauthClient)
		_ = client
	*/
	sshConfig := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			sshAgent(),
		},
	}

	connection, err := ssh.Dial("tcp", "192.241.137.221:22", sshConfig)
	if err != nil {
		log.Fatal(err)
	}

	session, err := connection.NewSession()
	if err != nil {
		log.Fatal(err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		session.Close()
		log.Fatal(err)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	go io.Copy(os.Stdout, stdout)

	err = session.Run("ls -ladh .")
	if err != nil {
		log.Fatal(err)
	}
	session.Close()
}

func sshAgent() ssh.AuthMethod {
	if sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		return ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers)
	}
	return nil
}
