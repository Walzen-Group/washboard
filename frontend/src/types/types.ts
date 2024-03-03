interface Container {
  id: string;
  name: string;
  image: string;
  upToDate: string;
  upToDateIgnored: boolean,
  status: string;
  ports: number[];
  labels: Record<string, string>;
}


interface Stack {
  id: number;
  name: string;
  containers: Container[];
  updateStatus: Object[];
}

interface ContainerView {

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

interface UpdateQueue extends Record<QueueStatus, Record<string, QueueItem>> { }

enum QueueStatus {
  Error = "error",
  Queued = "queued",
  Done = "done"
}

enum ImageStatus {
  Updated = "updated",
  Outdated = "outdated",
  Unavailable = "unavailable"
}

enum ContainerStatus {
  Updated = "updated",
  Outdated = "outdated",
  Skipped = "skipped",
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
  ContainerView,
  Snackbar,
  AppSettings,
  DockerUpdateManagerSettings,
  IgnoredImages,
  UpdateQueue, QueueItem,
  SidebarSettings
};
