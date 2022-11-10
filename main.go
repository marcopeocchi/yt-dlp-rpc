package main

import "goytdlp.rpc/m/pkg"

//	          __            ____         ____  ____  ______
//	   __  __/ /_      ____/ / /___     / __ \/ __ \/ ____/
//	  / / / / __/_____/ __  / / __ \   / /_/ / /_/ / /
//	 / /_/ / /_/_____/ /_/ / / /_/ /  / _, _/ ____/ /___
//	 \__, /\__/      \__,_/_/ .___/  /_/ |_/_/    \____/
//	/____/                 /_/
//
// Pluggable RPC for yt-dlp. Designed to handle many concurrent
// processes with little overhead as possible.
//
// https://github.com/marcopeocchi/yt-dlp-rpc
func main() {
	pkg.Run()
}
