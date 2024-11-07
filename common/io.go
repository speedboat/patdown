package common

import (
	"fmt"
	"os"
)

var (
	ColorReset     = "\033[0m"
	ColorRed       = "\033[31m"
	ColorPurple    = "\033[35m"
	ColorLightBlue = "\033[34m"
	ColorCyan      = "\033[36m"
	ColorGreen     = "\033[32m"
	ColorOrange    = "\033[91m"
	ColorGray      = "\033[90m"
	ColorYellow    = "\033[93m"
	ColorWhite     = "\033[97m"
)

func Usage() {
	Banner()
	fmt.Printf(`
usage: 
 %s!%s d | target fqdn (not recommended)
 %s!%s n | nameserver to query (can be specified multiple times)
   v | enable verbosity %s[false]%s
   t | threads %s[5]%s
   s | delay between requests in milliseconds, per thread %s[250]%s

e.g.
    patdown -d target.network
    patdown -n egress.ns.target.network -n another.egress.ns.target.network
    patdown -n dc.target.network -v -t 25
`, ColorRed, ColorReset, ColorRed, ColorReset, ColorPurple, ColorReset, ColorPurple, ColorReset, ColorPurple, ColorReset)
}

func Banner() {
	fmt.Fprintf(os.Stderr, `
                                  _______
             _/_   /          ---'   ____)____
    _   __.  /  __/ __ , , , ___        ______)
   /_)_(_/|_<__(_/_(_)(_(_/_/ <_        _______)
  /                                    _______)
 '                            ---.__________)

`)
}

func Success(msg string) {
	fmt.Printf("%s[+]%s %s\n", ColorGreen, ColorReset, msg)
}

func Info(msg string) {
	fmt.Printf("%s[i]%s %s\n", ColorCyan, ColorReset, msg)
}

func Warn(msg string) {
	fmt.Printf("%s[!]%s %s\n", ColorYellow, ColorReset, msg)
}

func Error(msg string) {
	fmt.Printf("%s[x]%s %s\n", ColorRed, ColorReset, msg)
}

func Fatal(msg string) {
	fmt.Printf("%s[f]%s %s\n", ColorRed, ColorReset, msg)
	os.Exit(-1)
}
