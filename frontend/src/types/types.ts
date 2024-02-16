interface Container {
  id: string;
  name: string;
  image: string;
  upToDate: string;
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
