package command

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/kemokemo/gckdir/lib"
	"github.com/skratchdot/open-golang/open"
	"github.com/urfave/cli"
)

var (
	// UsageVerify is Usage of verify subcommand for cli
	UsageVerify = "Verifies the structure and each hash value of files."
)

// CmdVerify verifies directory information below cases.
//  Case 1. a json file of hash list with target directory
//  Case 2. source directory with target directory
func CmdVerify(c *cli.Context) error {
	help := fmt.Sprintf("Please see '%s %s --help'.", c.App.Name, c.Command.FullName())
	source := c.Args().Get(0)
	target := c.Args().Get(1)
	if source == "" || target == "" {
		return cli.NewExitError(
			fmt.Sprintf("Source path or target path is empty. %s", help),
			ExitCodeInvalidArguments)
	}
	source = filepath.Clean(source)
	target = filepath.Clean(target)

	sourceList, err := lib.GetHashList(source)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("Failed to get the hash list. %v\n%s", err, help),
			ExitCodeFunctionError)
	}
	targetList, err := lib.GetHashList(target)
	if err != nil {
		return cli.NewExitError(
			fmt.Sprintf("Failed to get the hash list. %v\n%s", err, help),
			ExitCodeFunctionError)
	}

	result := lib.VerifyHashList(sourceList, targetList)
	var path string
	if c.Bool("report") || c.Bool("open") {
		pathList := lib.PathList{SourcePath: source, TargetPath: target}
		path, err = createReport(pathList, result)
		if err != nil {
			return cli.NewExitError(
				fmt.Sprintf("Failed to create a result report. %v\n%s", err, help),
				ExitCodeFunctionError)
		}
	}
	if c.Bool("open") {
		err = open.Run(path)
		if err != nil {
			return cli.NewExitError(
				fmt.Sprintf("Failed to open a result report. %v\n%s", err, help),
				ExitCodeFunctionError)
		}
	}

	if result.VerifyResult == false {
		log.Println("Verification failed.")
		return cli.NewExitError("", ExitCodeVerificationFailed)
	}
	return nil
}

func createReport(pathList lib.PathList, result lib.HashList) (string, error) {
	cd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	path := filepath.Join(cd, time.Now().Format("Result_20060102-030405.html"))

	file, err := os.Create(path)
	defer func() {
		err = file.Close()
		if err != nil {
			log.Println("Failed to close file", err)
		}
	}()

	err = lib.CreateReport(file, pathList, result)
	if err != nil {
		return "", err
	}

	path = filepath.Join("file:///", path)
	return path, nil
}
