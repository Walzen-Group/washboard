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
	Ports    []int                  `json:"ports"`
	Labels   map[string]interface{} `json:"labels"`
}

type StackDto struct {
	Id         int             `json:"id"`
	Name       string          `json:"name"`
	Containers []*ContainerDto `json:"containers"`
}

type StackUpdateStatus struct {
	EndpointId int    `json:"endpointId"`
	StackId    int    `json:"stackId"`
	StackName  string `json:"stackName"`
	Status     string `json:"status"`
	Details    string `json:"details"`
	Timestamp  int64  `json:"timestamp"`
}

type ActionType string

// We could make an ActionType type and use that instead of string but that would require some annoying refactoring
const (
	Outdated                  string     = "outdated"
	Updated                   string     = "updated"
	Preparing                 string     = "preparing"
	Skipped                   string     = "skipped"
	Error                     string     = "error"
	Done                      string     = "done"
	Queued                    string     = "queued"
	DbName                    string     = "washb"
	DbGroupSettingsCollection string     = "group_settings"
	DbStackSettingsCollection string     = "stack_settings"
	Start                     ActionType = "start"
	Stop                      ActionType = "stop"
	Kill                      ActionType = "kill"
	Restart                   ActionType = "restart"
	Pause                     ActionType = "pause"
	Resume                    ActionType = "resume"
)

type GroupSettings struct {
	GroupName      string `bson:"groupName" json:"groupName"`
	GlobalPriority int    `bson:"globalPriority" json:"globalPriority"`
}

type StackSettings struct {
	StackId             int `bson:"stackId" json:"stackId"`
	GlobalPriority      int `bson:"globalPriority" json:"globalPriority"`
	WithinGroupPriority int `bson:"withinGroupPriority" json:"withinGroupPriority"`
}
