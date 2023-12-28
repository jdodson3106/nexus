package nexus

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	CREATE_MARKER      = color.HiGreenString(" +")
	DELETE_MARKER      = color.HiRedString(" -")
	INFO_CARET         = color.CyanString(" >")
	NOTE_ALERT         = color.HiYellowString(" NOTE:")
	WARNING_INFO_ALERT = color.HiRedString(" !")
	WARNING_MARKER     = color.HiYellowString(" ?")
)

func PrintInfo(cmd string) {
	fmt.Printf("%s %s", INFO_CARET, cmd)
}

func PrintWarningInfo(cmd string) {
	fmt.Printf("%s %s", WARNING_INFO_ALERT, cmd)
}

func PrintCreate(cmd string) {
	fmt.Printf("%s %s", CREATE_MARKER, cmd)
}

func PrintDelete(cmd string) {
	fmt.Printf("%s %s", DELETE_MARKER, cmd)
}

func PrintNote(cmd string) {
	fmt.Printf("%s %s", NOTE_ALERT, cmd)
}

func PrintWarning(cmd string) {
	fmt.Printf("%s %s", WARNING_MARKER, cmd)
}
