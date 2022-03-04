package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	flag2 "github.com/CarsonSlovoka/fox/pkg/flag"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var (
	ChangeDirError = errors.New("change dir error")
)

func initCDCmd() *flag2.Command {
	cmd := flag2.NewCommand(flag.NewFlagSet("cd", flag.ContinueOnError),
		map[string][]flag2.CmdField{
			"string": {
				{"wDir", "", "Working Directory"},
			},
		})
	cmd.MainFunc = func(args []string) error {
		if err := cmd.Parse(args, true); err != nil {
			return err
		}

		wDir := (cmd.Lookup("wDir")).Value.(flag.Getter).Get().(string)
		if wDir != "" {
			if err := os.Chdir(wDir); err != nil {
				return ChangeDirError
			}
			workingDir, err := os.Getwd()
			if err != nil {
				return err
			}
			fmt.Println(fmt.Sprintf("Working Directory:%s", workingDir))
		}
		return nil
	}
	return cmd
}

func initGoCmd() *flag2.Command {
	cmd := flag2.NewCommand(flag.NewFlagSet("go", flag.ContinueOnError),
		map[string][]flag2.CmdField{
			"string": {
				{"wDir", "", "Working Directory"},
			},
			"int": {
				{"utils", 0, "select run case"},
			},
		})
	cmd.Usage = func() {
		_, _ = fmt.Fprintf(cmd.FlagSet.Output(),
			"%s: Usage of %s:\n"+"",
			cmd.FlagSet.Name(),
		)
		cmd.FlagSet.PrintDefaults()
	}
	cmd.MainFunc = func(args []string) error {
		if err := cmd.Parse(args, true); err != nil {
			return err
		}

		wDir := (cmd.Lookup("wDir")).Value.(flag.Getter).Get().(string)
		if wDir != "" {
			if err := os.Chdir(wDir); err != nil {
				return ChangeDirError
			}
			workingDir, err := os.Getwd()
			if err != nil {
				return err
			}
			fmt.Println(fmt.Sprintf("Working Directory:%s", workingDir))
		}
		return nil
	}
	return cmd
}

func startCMD(quitChan *chan error) {
	cmdCD := initCDCmd()
	var flagSetAll []*flag.FlagSet
	for _, curCmd := range []*flag2.Command{cmdCD} {
		flagSetAll = append(flagSetAll, curCmd.FlagSet)
	}

	menuHelp := func(args []string) error {
		for _, curFlagSet := range flagSetAll {
			curFlagSet.Usage()
		}
		return nil
	}

	type msgFunc func(args []string) error
	msgMap := map[string]msgFunc{
		"help":  menuHelp,
		"-help": menuHelp,
		"-h":    menuHelp,
		"quit": func(args []string) error {
			*quitChan <- errors.New("terminal close")
			return nil
		},
		"cls": func(args []string) error {
			var clearMap map[string]func() error
			clearMap = make(map[string]func() error)
			clearMap["linux"] = func() error {
				cmd := exec.Command("clear")
				cmd.Stdout = os.Stdout
				return cmd.Run()
			}
			clearMap["windows"] = func() error {
				cmd := exec.Command("cmd", "/c", "cls") // /c: Close
				cmd.Stdout = os.Stdout
				return cmd.Run()
			}
			clearFunc, ok := clearMap[runtime.GOOS]
			if !ok {
				return errors.New("your platform is unsupported! i can't clear terminal screen :(")

			}
			return clearFunc()
		},
	}

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Enter CMD: ")
		scanner.Scan()
		text := scanner.Text()
		args := strings.Split(text, " ")
		if handleFunc, exists := msgMap[args[0]]; exists {
			if err := handleFunc(args[1:]); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err.Error())
			}
		}
	}
}
