package main

import (
	"fmt"
	"github.com/michele-mogul/mutasync/internal/mutasync"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

import (
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting Mutagen setup, have gun! (EHM FUN)")
	fmt.Println("Shutdown docker and docker-sync")

	dockerComposePath := askForPath("docker compose")
	dockerSyncPath := askForPath("docker sync")
	envFile := askForPath(".env")
	dir, _ := filepath.Abs(filepath.Dir(envFile))

	command := exec.Command("mutasync", "sync", "list")
	out, _ := command.Output()
	mutagenList := string(out)
	r := regexp.MustCompile(`sync_(.*)`)
	result := r.FindAllString(mutagenList, -1)
	if len(result) > 0 {
		for _, mutaSync := range result {
			exec.Command("mutasync",
				"sync",
				"terminate",
				mutaSync).Run()
		}
	}

	exec.Command("docker-compose",
		"-f",
		dockerComposePath,
		"--project-directory",
		dir,
		"down").Run()

	exec.Command("docker-sync",
		"-c",
		dockerSyncPath,
		"stop").Run()

	exec.Command("docker-compose",
		"-f",
		dockerComposePath,
		"--project-directory",
		dir,
		"up",
		"-d").Run()

	_ = godotenv.Load()
	var myEnv = make(map[string]string)
	var syncs = make(map[string]string)

	myEnv, _ = godotenv.Read(envFile)
	projectRoot := myEnv["REPOPATH"]
	fmt.Println(projectRoot)
	compose, errCompose := mutasync.ParseCompose(dockerComposePath)
	sync, errSync := mutasync.ParseSync(dockerSyncPath)
	fmt.Println(sync)
	if errCompose != nil || errSync != nil {
		fmt.Println(errCompose)
		fmt.Println(errSync)
	} else {

		for k, value := range sync.Syncs {
			syncs[k] = value.Src
		}
		//Mutagen Commands generate
		generateMutagenCommands(compose, syncs, myEnv)

	}
}

func generateMutagenCommands(compose *mutasync.Compose, syncs map[string]string, myEnv map[string]string) {
	for _, value := range compose.Services {
		fmt.Println(value)
		volumes := value.Volumes
		for _, valueVolume := range volumes {
			switch valueVolume.(type) {
			case string:
				str := fmt.Sprintf("%v", valueVolume)
				var split = strings.Split(str, ":")
				if _, ok := syncs[split[0]]; ok {
					workingDir := value.WorkingDir
					if len(workingDir) == 0 {
						workingDir = split[1]
					}
					realPath := syncs[split[0]]
					//Removing env variables
					for keyReplace, replace := range myEnv {
						workingDir = strings.Replace(workingDir, "${"+keyReplace+"}", replace, -1)
						realPath = strings.Replace(realPath, "${"+keyReplace+"}", replace, -1)
					}
					localMountPath := strings.Replace(workingDir, split[1], realPath, -1)
					remoteMountPath := strings.Replace(workingDir, realPath, split[1], -1)
					mutagenSyncIdentifier := value.ContainerName + "-mutagensync"
					mutagenCommand := mutasync.MutagenCommand{
						Name:           mutagenSyncIdentifier,
						Ignore:         nil,
						ContainerName:  value.ContainerName,
						ContainerPath:  remoteMountPath,
						LocalMountPath: localMountPath,
					}
					//TODO link to the right main container
					errorMutagen := mutasync.CreateCommand(mutagenCommand)
					fmt.Println(errorMutagen)
				}
				break
			}
		}
	}
}

func askForPath(labelPath string) string {

	var (
		condition         bool
		pathRequiredInput string
	)
	condition = false

	for ok := true; ok; ok = condition {
		fmt.Println("Enter " + labelPath + " path")

		fmt.Scanln(&pathRequiredInput)
		condition = exits(pathRequiredInput)
		if condition {
			fmt.Println("Path not valid, enter correct path")
		}
	}

	return pathRequiredInput
}
