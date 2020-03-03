package mutasync


type Compose struct {
	Version  string
	Networks map[string]Network
	Volumes  map[string]Volume
	Services map[string]Service
}

type Network struct {
	Driver, External string
	DriverOpts       map[string]string "driver_opts"
}

type Volume struct {
	Driver, External string
	DriverOpts       map[string]string "driver_opts"
}

type Service struct {
	ContainerName                     string "container_name"
	Volumes			[]interface{}     "volumes"
	WorkingDir		string "working_dir"
}

type Sync struct {
	Version  string
	Syncs map[string]Syncs
}

type Syncs struct {
	SyncStrategy	string   	 "sync_strategy"
	Src				string  	 "internal"
	Exclude			[]string	 "sync_excludes"
}

type MutagenCommand struct {
	Name           string
	Ignore         []string
	User           string
	ContainerName  string
	ContainerPath  string
	LocalMountPath string
}