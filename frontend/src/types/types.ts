interface Container {
  id: string;
  name: string;
  image: string;
  upToDate: string;
  upToDateIgnored: boolean,
  status: string;
  ports: number[];
  labels: { [key: string]: string };
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
  id: number;
  message: string;
  color: string;
  show: boolean;
}

interface AppSettings {
  dockerUpdateManagerSettings: DockerUpdateManagerSettings;
}

interface IgnoredImages {
  [key: string]: boolean;
}

interface DockerUpdateManagerSettings {
  ignoredImages: IgnoredImages;
}

interface UpdateStackQueue {
  [key: string]: {
    Expiration: number;
    Object: {
      details: string;
      status: string;
      endpointId: number;
      stackId: number;
    }
  };
}
