package gormmigrations

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var args *cmdArgs
var flagsValid = []string{"name", "table", "path"}

type cmdArgs struct {
	command string
	params  map[string]string
	flags   map[string]string
}

func setUpArgs() *cmdArgs {
	params, flags := setArgs()
	return &cmdArgs{
		command: flag.Arg(0),
		params:  params,
		flags:   flags,
	}
}

func setArgs() (map[string]string, map[string]string) {
	params := make(map[string]string)
	flags := make(map[string]string)
	switch flag.Arg(0) {
	case "migrations:migrate":
		break
	case "migrations:rollback":
		break
	case "migrations:rollback:last":
		break
	case "migrations:create":
		flagsCmd := setFlags(os.Args, 3)
		if flag.Arg(1) == "" || strings.HasPrefix(flag.Arg(1), "--") {
			log.Fatal("Name of migration is required")
		}
		params["migration"] = flag.Arg(1)

		table := ""
		value, exists := flagsCmd["table"]
		if exists {
			table = value
		}

		path := ""
		value, exists = flagsCmd["path"]
		if exists {
			path = value
		}
		flags["table"] = table
		flags["path"] = path
	case "migrations:seeder":
	default:
		log.Fatalf("Command %s unknow", args.command)
	}
	return params, flags
}

func setFlags(args []string, pos int) map[string]string {
	var params = make(map[string]string)
	if len(args) > pos {
		for i := pos; i < len(args); i++ {
			valid, err := validateFlag(args[i])
			if err != nil {
				fmt.Println(err.Error())
			}
			if valid {
				value := strings.Split(strings.TrimPrefix(args[i], "--"), "=")
				params[value[0]] = value[1]
			}
		}
	}

	return params
}

func validateFlag(param string) (bool, error) {
	if !strings.HasPrefix(param, "--") {
		return false, errors.New("Unknow flag " + param)
	}

	if !strings.Contains(param, "=") {
		err := fmt.Sprintf("The param %s must be have a value", param)
		return false, errors.New(err)
	}

	value := strings.Split(strings.TrimPrefix(param, "--"), "=")
	exists := false

	for _, v := range flagsValid {
		if value[0] == v {
			exists = true
			break
		}
	}

	if !exists {
		err := fmt.Sprintf("Unknow parameter %s", param)
		return false, errors.New(err)
	}

	if value[1] == "" {
		err := fmt.Sprintf("The param %s must be have a valid value", param)
		return false, errors.New(err)
	}

	return true, nil
}
