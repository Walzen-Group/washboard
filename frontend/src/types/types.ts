interface Container {
  id: string;
  name: string;
  image: string;
  upToDate: boolean;
  status: string;
  ports: number[];
  labels: { [key: string]: string };
}

interface Stack {
  id: number;
  name: string;
  containers: Container[];
}
