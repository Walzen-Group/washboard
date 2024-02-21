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

interface UpdateQueue extends Record<QueueStatus, Record<string, QueueItem>> { }

enum QueueStatus {
  error = "error",
  queued = "queued",
  done = "done"
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
  QueueStatus
};
export type {
  Container,
  Stack,
  ContainerView,
  Snackbar,
  AppSettings,
  DockerUpdateManagerSettings,
  IgnoredImages,
  UpdateQueue, QueueItem
};
