package nexus

import (
	"fmt"
	"os/exec"
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
}

func installTempl() error {
    const templLib = "github.com/a-h/templ"
    out, err := exec.Command("go", "version").Output()
    if err != nil {
        panic(err)
    }

    fmt.Printf("Running with Go version %s\n", out)
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
            fmt.Printf("%s already installed.\n", dep)
            break
        }
    }

    if !hasTempl {
        fmt.Printf("installing %s\n", templLib)
        err = exec.Command("go", "install", "github.com/a-h/templ/cmd/templ@latest").Run()
        if err != nil {
            fmt.Printf("Error installing %s :: %s\n", templLib, err)
            return err
        }

        err = exec.Command("go", "mod", "tidy").Run()
        if err != nil {
            fmt.Printf("Error runing mod tidy :: %s\n", err)
            return err
        }
    }
    
    return nil
}

func compileTemplates() error {
    return exec.Command("templ", "generate").Run()
}
