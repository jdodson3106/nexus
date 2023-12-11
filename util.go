package nexus

import (
	"fmt"
	"os/exec"

	"github.com/jdodson3106/nexus/log"
)

func printAppString() {
	name := ` 
 ________    _______       ___    ___  ___  ___   ________      
|\   ___  \ |\  ___ \     |\  \  /  /||\  \|\  \ |\   ____\     
\ \  \\ \  \\ \   __/|    \ \  \/  / /\ \  \\\  \\ \  \___|_    
 \ \  \\ \  \\ \  \_|/__   \ \    / /  \ \  \\\  \\ \_____  \   
  \ \  \\ \  \\ \  \_|\ \   /     \/    \ \  \\\  \\|____|\  \  
   \ \__\\ \__\\ \_______\ /  /\   \     \ \_______\ ____\_\  \ 
    \|__| \|__| \|_______|/__/ /\ __\     \|_______||\_________\
                          |__|/ \|__|               \|_________|
	`
	fmt.Printf("\033[1;32m%s\033[0m", name)
    fmt.Println()
}

func tidy() error {
    err := exec.Command("go", "mod", "tidy").Run()
    if err != nil {
        log.Error(fmt.Sprintf("Error runing mod tidy :: %s", err))
        return err
    }

    return nil
}

func installTempl() error {
    const templLib = "github.com/a-h/templ"
    out, err := exec.Command("go", "version").Output()
    if err != nil {
        panic(err)
    }

    log.Trace(fmt.Sprintf("Running with Go version %s", out))
    out, err = exec.Command("go", "list", "-deps").Output()
    depBuffer := []byte{}
    deps := []string {}
    for _, v := range out {
        if v == '\n' {
            deps = append(deps, string(depBuffer))
            depBuffer = nil
            continue
        }
        depBuffer = append(depBuffer, v)
    }

    hasTempl := false 
    for _, dep := range deps {
        if dep == templLib {
            hasTempl = true
            log.Trace(fmt.Sprintf("%s already installed.", dep))
            break
        }
    }

    if !hasTempl {
        log.Trace(fmt.Sprintf("installing %s", templLib))
        err = exec.Command("go", "install", "github.com/a-h/templ/cmd/templ@latest").Run()
        if err != nil {
            log.Error(fmt.Sprintf("Error installing %s :: %s", templLib, err))
            return err
        }
    }
    tidy()
    return nil
}

func compileTemplates() error {
    return exec.Command("templ", "generate").Run()
}
