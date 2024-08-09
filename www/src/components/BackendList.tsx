import { IServer } from "../types";

export default function BackendList({ data }: { data: Array<IServer> }) {
  return (
    <div className="border rounded-md overflow-hidden">
      {/* <p className="text-sm font-semibold">Backend List</p> */}
      <table className="table-auto w-full">
        <thead className="bg-slate-500">
          <tr className="text-sm">
            <th className="text-left p-3">Endpoint</th>
            <th className="py-3">Health</th>
            <th className="p-3">Status</th>
            <th className="p-3">Weight</th>
            <th className="p-3">Active Requests</th>
            <th className="p-3">Success Rate</th>
            <th className="p-3">Error Rate</th>
          </tr>
        </thead>
        <tbody>
          {data.map((server) => (
            <tr className="text-sm font-medium">
              <td className="p-3">{server.url}</td>
              <td className="text-center">Healthy</td>
              <td className="text-center">Status</td>
              <td className="text-center">{server.weight}</td>
              <td className="text-center">{server.activeConnections}</td>
              <td className="text-center">0</td>
              <td className="text-center">0</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
