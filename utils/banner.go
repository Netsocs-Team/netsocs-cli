package utils

import "github.com/pterm/pterm"

func ShowBannerArt() {
	art := ` 
███╗   ██╗███████╗████████╗███████╗ ██████╗  ██████╗███████╗
████╗  ██║██╔════╝╚══██╔══╝██╔════╝██╔═══██╗██╔════╝██╔════╝
██╔██╗ ██║█████╗     ██║   ███████╗██║   ██║██║     ███████╗
██║╚██╗██║██╔══╝     ██║   ╚════██║██║   ██║██║     ╚════██║
██║ ╚████║███████╗   ██║   ███████║╚██████╔╝╚██████╗███████║
╚═╝  ╚═══╝╚══════╝   ╚═╝   ╚══════╝ ╚═════╝  ╚═════╝╚══════╝
`
	pterm.DefaultCenter.Println(pterm.NewRGB(0, 128, 0).Sprint("NETSOCS"))
	pterm.DefaultCenter.Println(pterm.NewRGB(0, 128, 0).Sprint("Server Configuration System"))
	pterm.DefaultCenter.Println(pterm.NewRGB(0, 128, 0).Sprint("Versión 1.0"))
	pterm.DefaultCenter.Println(pterm.NewRGB(0, 128, 0).Sprint(""))
	pterm.DefaultCenter.Println(art)
}
