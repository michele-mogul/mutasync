package mutasync

import (
	"fmt"
	"os/exec"
	"strings"
)



func CreateCommand(command MutagenCommand) error {

	var ignore string
	if len(command.Ignore) > 0 {
		ignore =  "-i " + strings.Join(command.Ignore, " ")
	}
	dockerCommand := "docker://"+command.ContainerName+command.ContainerPath+" "+command.LocalMountPath

	cmd := exec.Command(
		"typos",
		"sync",
		"create",
		"--name="+command.Name,
		"-m=two-way-resolved",
		"--symlink-mode ignore",
		ignore,
		dockerCommand)
	err := cmd.Run()
	fmt.Println(cmd.Stdout)
	return err

}

func buildMutagenCommand(user string, container string, containerpath string, local string, name string, ignore []string) *MutagenCommand {

	mutagenCommand := &MutagenCommand{name,ignore, user,container,containerpath,local}
	return mutagenCommand
}