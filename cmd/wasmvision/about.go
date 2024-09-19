package main

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

var logo = `

                             ____   ____.__       .__               
__  _  _______    ______ ____\   \ /   /|__| _____|__| ____   ____  
\ \/ \/ /\__  \  /  ___//     \   Y   / |  |/  ___/  |/  _ \ /    \ 
 \     /  / __ \_\___ \|  Y Y  \     /  |  |\___ \|  (  <_> )   |  \
  \/\_/  (____  /____  >__|_|  /\___/   |__/____  >__|\____/|___|  /
              \/     \/      \/                 \/               \/ 

wasmVision - gets you up and running with computer vision using WebAssembly.

https://wasmvision.com
`

func about(cCtx *cli.Context) error {
	fmt.Println(logo)
	fmt.Println("Version:", Version())

	return nil
}
