package types

type EndpointDto struct {
	Id         int            `json:"id"`
	Name       string         `json:"name"`
	Containers []ContainerDto `json:"containers"`
}

type GenericDto struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type ContainerDto struct {
	Id       string                 `json:"id"`
	Name     string                 `json:"name"`
	Image    string                 `json:"image"`
	UpToDate string                 `json:"upToDate"`
	Status   string                 `json:"status"`
	Networks []string               `json:"networks"`
	Ports    []int                  `json:"ports"`
	Labels   map[string]interface{} `json:"labels"`
}

type StackDto struct {
	Id         int             `json:"id"`
	Name       string          `json:"name"`
	Containers []*ContainerDto `json:"containers"`
	Priority   int             `json:"priority"`
	AutoStart  bool            `json:"autoStart"`
}

type StackUpdateStatus struct {
	EndpointId int    `json:"endpointId"`
	StackId    int    `json:"stackId"`
	StackName  string `json:"stackName"`
	Status     string `json:"status"`
	Details    string `json:"details"`
	Timestamp  int64  `json:"timestamp"`
}

type ContainerAction string

// We could make an ActionType type and use that instead of string but that would require some annoying refactoring
const (
	Outdated                  string          = "outdated"
	Updated                   string          = "updated"
	Preparing                 string          = "preparing"
	Skipped                   string          = "skipped"
	Error                     string          = "error"
	Done                      string          = "done"
	Queued                    string          = "queued"
	NotRequested              string          = "not_requested"
	DbName                    string          = "washb"
	DbGroupSettingsCollection string          = "group_settings"
	DbStackSettingsCollection string          = "stack_settings"
	DbAccountsCollection      string          = "accounts"
	StackGroupLabel           string          = "org.walzen.washb.webui"
	WebUIMachineAddressKey    string          = "${ADDRESS}"
	StackLabel                string          = "com.docker.compose.project"
	IdentityKey               string          = "id"
	Start                     ContainerAction = "start"
	Stop                      ContainerAction = "stop"
	Kill                      ContainerAction = "kill"
	Restart                   ContainerAction = "restart"
	Pause                     ContainerAction = "pause"
	Resume                    ContainerAction = "resume"
)

type Login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

type User struct {
	UserName string
}
type StackSettings struct {
	StackName string `bson:"stackName" json:"stackName"`
	Priority  int    `bson:"globalPriority" json:"globalPriority"`
	AutoStart bool   `bson:"autoStart" json:"autoStart"`
}

type SyncOptions struct {
	EndpointIds []int `json:"endpointIds"`
}
