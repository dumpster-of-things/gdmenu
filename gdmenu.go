package gdmenu

// Source:   git clone https://git.suckless.org/dmenu
// Patches: grid(myFixed), gridnav, multi-select

import (
	"io"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type Grid struct {
	x int //columns
	y int //rows
}
func (i Grid) cols() string {
	return strconv.Itoa(i.x)
}
func (i Grid) rows() string {
	return strconv.Itoa(i.y)
}

func MkPrompt(prompt string) ([]byte, error) {
	// initialize dmenu prompt
	cmd := exec.Command("dmenu", "-p", prompt)
	stdin, err := cmd.StdinPipe()
	if err != nil { log.Fatal(err) }
	// pass opts to dmenu
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "")
	}()
	// return stdout, stderr
	return(cmd.CombinedOutput())
}

func Prompt(text string) string {
	var Out string
	output, err := MkPrompt(text)
	if err != nil || string(output) == "\n" { Out = "null" } else { Out = []string(strings.Split(string(output), "\n"))[0] }
	return Out
}

func MkMenu(options []string) ([]byte, error) {
	// initialize grid dimensions
	grid := Grid{}
	maxcols := 3 //max columns
	maxrows := 9 // menu legnth (rows) limit
	n := len(options) // total menu options (possible rows)

	if a := (n / maxcols); a >= maxrows {
		grid.x = maxcols
		grid.y = maxrows
	} else {
		if n % maxcols == 0 {
			grid.x = maxcols
			grid.y = a
		} else {
			if mp := (maxcols + 1); n % mp == 0 {
				grid.x = mp
				grid.y = a
			} else {
				grid.x = maxcols
				grid.y = a+1
			}
		}
	}

	// initialize dmenu with grid dimensions
	cmd := exec.Command("dmenu", "-i", "-g", grid.cols(), "-l", grid.rows())
	stdin, err := cmd.StdinPipe()
	if err != nil { log.Fatal(err) }

	// pass opts to dmenu
	go func() {
		defer stdin.Close()
		io.WriteString(stdin, strings.Join(options, "\n"))
	}()

	// return stdout, stderr
	return(cmd.CombinedOutput())
}

func SafeSelect(Opts []string) []string {
	var Outputs []string
	out, err := MkMenu(Opts)
	if err != nil || string(out) == "\n" {
		Outputs = []string {"null"}
	} else {
		// Convert "out"([]byte) to "Output"([]string)
		//Outputs = []string(strings.Split(string(out), "\n"))
		for _, o := range []string(strings.Split(string(out), "\n")) {
			switch o {
			case "Cancel":
				Outputs = []string {"null"}
				break
			case "Done":
				continue
			case "\n":
				continue
			case "":
				continue
			default:
				Outputs = append(Outputs, o)
				break
			}
		}
	}
	//return Outputs[:(len(Outputs)-1)]
	return Outputs
}

func SafeSwitch(Opts []string) string {
	var Output string
	out, err := MkMenu(Opts)
	if err != nil || string(out) == "\n" {
		Output = "null"
	} else {
		// Convert "out"([]byte) to "Output"([]string)
		//Output = []string(strings.Split(string(out), "\n"))
		//for _, o := range SafeSelect(Opts) {
		for _, o := range []string(strings.Split(string(out), "\n")) {
			switch o {
			case "Cancel":
				Output = "null"
				break
			case "Done":
				continue
			case "\n":
				continue
			case "":
				continue
			default:
				Output = o
				break
			}
		}
	}
	return Output

}

/*func (m *Menu) SafeMenu() {

}*/

func SelectFromSlice(Opts []string) ([]byte, error) { return(MkMenu(Opts)) }

func SwitchFromSlice(Opts []string) string { return(SafeSwitch(Opts)) }

func SelectFromList(Opts ...string) ([]byte, error) { return(MkMenu(Opts)) }

func SwitchFromList(Opts ...string) string { return(SafeSwitch(Opts)) }

