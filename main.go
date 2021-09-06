package main

import (
	"fmt"

	"github.com/levidurfee/peerr/pkg"
)

func main() {
	me := pkg.Start("./config.json")
	me.AddPeer()
	// fmt.Println(me, me.Peer)
	me.Wireguard()
	me.Bird()

	fmt.Println("#########################")
	fmt.Println()
	fmt.Printf("wg-quick up %s\n", me.Peer.Name)
	fmt.Printf("systemctl enable wg-quick@%s\n", me.Peer.Name)
	fmt.Println("birdc config check")
	fmt.Println("birdc config")
	fmt.Println()
	fmt.Println("#########################")
}
