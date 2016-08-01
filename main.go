package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"github.com/digitalocean/godo"

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

var token = flag.String("token", "", "Digital Ocean API Token")

func listDroplets(client *godo.Client) ([]godo.Droplet, error) {
	var list []godo.Droplet

	opt := &godo.ListOptions{}
	for {
		droplets, resp, err := client.Droplets.List(opt)
		if err != nil {
			return nil, err
		}

		list = append(list, droplets...)

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		opt.Page = page + 1
	}

	return list, nil
}

func main() {
	flag.Parse()
	tokenSource := AccessToken(*token)
	oauthClient := oauth2.NewClient(oauth2.NoContext, &tokenSource)
	client := godo.NewClient(oauthClient)

	droplets, err := listDroplets(client)
	if err != nil {
		log.Fatal(err)
	}

	var networks []string
	for _, d := range droplets {
		for _, network := range d.Networks.V4 {
			networks = append(networks, network.IPAddress)
		}
	}

	if len(networks) < 1 {
		log.Fatal("Can't ssh because there are no droplets under the specificed token.")
	}

	// for now, just get the first one...
	dropletIP := fmt.Sprintf("%s:22", networks[0])

	sshConfig := &ssh.ClientConfig{
		User: "root",
		Auth: []ssh.AuthMethod{
			sshAgent(),
		},
	}

	connection, err := ssh.Dial("tcp", dropletIP, sshConfig)
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
