export type IServer = {
  activeConnections: number;
  reqServedCount: number;
  url: string;
  weight: number;
};
