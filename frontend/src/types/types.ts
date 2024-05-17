interface Container {
  id: string;
  name: string;
  image: string;
  upToDate: string;
  upToDateIgnored: boolean,
  status: string;
  networks: string[];
  ports: string[];
  labels: Record<string, string>;
}

interface StackInternal extends Stack {
  expanded: boolean;
  checked: boolean;
}

interface Stack extends StackSettings {
  id: number;
  name: string;
  containers: Container[];
  updateStatus: Object[];
}

interface StackSettings {
  priority: number;
  autoStart: boolean;
}

interface StackSettingsDto {
  stackName: string;
  stackId: number;
  priority: number;
  autoStart: boolean;
}

interface Group extends GroupSettings {
  stacks: Stack[];
}

interface GroupSettings {
  groupName: string;
}

interface PaddingClass {
  [key: string]: string | undefined;
}

interface Snackbar {
  id: string;
  message: string;
  color: string;
  show: boolean;
}

interface AppSettings {
  dockerUpdateManagerSettings: DockerUpdateManagerSettings;
}

interface IgnoredImages extends Record<string, boolean> { }

interface DockerUpdateManagerSettings {
  ignoredImages: IgnoredImages;
}

interface SidebarSettings {
  show?: boolean;
  mini?: boolean;
  clipped?: boolean;
}

interface URLConfig {
  defaultHost: string;
  defaultPortainerAddress: string;
}

interface UpdateQueue extends Record<QueueStatus, Record<string, QueueItem>> { }

enum QueueStatus {
  Error = "error",
  Queued = "queued",
  Done = "done"
}

enum ImageStatus {
  Updated = "updated",
  Outdated = "outdated",
  NotRequested = "not_requested",
  Unavailable = ""
}

enum ContainerStatus {
  Updated = "updated",
  Outdated = "outdated",
  Skipped = "skipped",
  NotRequested = "not_requested",
  Error = "error"
}

enum Action {
  Start = "start",
  Stop = "stop",
  Restart = "restart",
}

interface QueueItem {
  details: string;
  status: QueueStatus;
  stackName: string;
  endpointId: number;
  stackId: number;
  timestamp: number;
}

export {
  QueueStatus,
  ImageStatus,
  ContainerStatus,
  Action
};
export type {
  Container,
  Stack,
  StackSettingsDto,
  StackInternal,
  StackSettings,
  GroupSettings,
  Group,
  PaddingClass,
  Snackbar,
  AppSettings,
  DockerUpdateManagerSettings,
  IgnoredImages,
  UpdateQueue, QueueItem,
  SidebarSettings,
  URLConfig
};
