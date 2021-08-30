package pkg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

// Node represents a node
type Node struct {
	ASN      string `json:"asn"`
	Name     string `json:"name"`
	Endpoint string `json:"endpoint"`
	Port     string `json:"port"`

	Comment string `json:"comment"`

	PublicIPs IPs `json:"public_ips"`
	DN42IPs   IPs `json:"dn42_ips"`

	WG WG `json:"wg"`

	Peer *Node `json:"peer"`

	Output Output `json:"output"`
}

type Output struct {
	WG   string `json:"wg"`
	Bird string `json:"bird"`
}

// Start _
func Start(file string) Node {
	var node Node
	jsonFile, err := os.Open(file)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Successfully Opened %s\n", file)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	bytes, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(bytes, &node)

	return node
}

// Create _
func Create() Node {
	var node Node

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Peer'r")
	fmt.Println("---------------------")

	fmt.Print("ASN -> ")
	text, _ := reader.ReadString('\n')
	node.ASN = parse(text)

	fmt.Print("Abbreviation-Region -> ")
	text, _ = reader.ReadString('\n')
	node.Name = parse(text)

	fmt.Print("Comment -> ")
	text, _ = reader.ReadString('\n')
	node.Comment = parse(text)

	fmt.Print("Remote Endpoint (include port) -> ")
	text, _ = reader.ReadString('\n')
	node.Endpoint = parse(text)

	fmt.Print("Wireguard IPv4 -> ")
	text, _ = reader.ReadString('\n')
	node.DN42IPs.V4 = parse(text)

	fmt.Print("Wireguard IPv6 -> ")
	text, _ = reader.ReadString('\n')
	node.DN42IPs.V6 = parse(text)

	fmt.Print("Public Key -> ")
	text, _ = reader.ReadString('\n')
	node.WG.PublicKey = parse(text)

	return node
}

func (n *Node) AddPeer() {
	peer := Create()
	n.Peer = &peer
}

// Wireguard _
func (n *Node) Wireguard() {
	f, err := os.Open("./templates/wireguard.conf")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	bytes, _ := ioutil.ReadAll(f)

	port := n.Peer.ASN[len(n.Peer.ASN)-4:]

	wgcfg := string(bytes)
	wgcfg = strings.Replace(wgcfg, "{{private_key}}", n.WG.PrivateKey, 1)
	wgcfg = strings.Replace(wgcfg, "{{dn42_ipv4}}", n.DN42IPs.V4, 1)
	wgcfg = strings.Replace(wgcfg, "{{dn42_ipv6}}", n.DN42IPs.V6, 1)
	wgcfg = strings.Replace(wgcfg, "{{peer_ipv4}}", n.Peer.DN42IPs.V4, 1)
	wgcfg = strings.Replace(wgcfg, "{{peer_ipv6}}", n.Peer.DN42IPs.V6, 1)
	wgcfg = strings.Replace(wgcfg, "{{listen_port}}", port, 1)
	wgcfg = strings.Replace(wgcfg, "{{peer_endpoint}}", n.Peer.Endpoint, 1)
	wgcfg = strings.Replace(wgcfg, "{{peer_public_key}}", n.Peer.WG.PublicKey, 1)
	wgcfg = strings.Replace(wgcfg, "{{comment}}", n.Peer.Comment, 1)

	// fmt.Println(wgcfg)
	// Check if file exists before writing to it. Don't overwrite files.
	wg, _ := os.Create(n.Output.WG + "/" + n.Peer.Name + ".conf")
	defer wg.Close()

	wg.WriteString(wgcfg)
}

// Bird _
func (n *Node) Bird() {
	f, err := os.Open("./templates/bird.conf")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	bytes, _ := ioutil.ReadAll(f)

	peerName := strings.Replace(n.Peer.Name, "-", "_", 1)

	cfg := string(bytes)
	cfg = strings.Replace(cfg, "{{abbreviation}}", peerName, 1)
	cfg = strings.Replace(cfg, "{{ip}}", n.Peer.DN42IPs.V6, 1)
	cfg = strings.Replace(cfg, "{{asn}}", n.Peer.ASN, 1)
	cfg = strings.Replace(cfg, "{{interface}}", n.Peer.Name, 1)

	brd, _ := os.Create(n.Output.Bird + "/" + n.Peer.Name + ".conf")
	defer brd.Close()

	brd.WriteString(cfg)
}

func parse(input string) string {
	// convert CRLF to LF
	return strings.Replace(input, "\n", "", -1)
}
