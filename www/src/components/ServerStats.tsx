import { useEffect, useState } from "react";

const ServerStats = () => {
  const [serverData, setServerData] = useState({});
  const [isConnected, setIsConnected] = useState(false);

  useEffect(() => {
    // Create a new WebSocket connection
    const ws = new WebSocket("ws://localhost:8080/ws");

    ws.onopen = () => {
      console.log("Connected to the WebSocket server");
      setIsConnected(true);
    };

    ws.onmessage = (event) => {
      const data = JSON.parse(event.data);
      setServerData(data);
    };

    ws.onclose = () => {
      console.log("Disconnected from the WebSocket server");
      setIsConnected(false);
    };

    ws.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    return () => {
      ws.close();
    };
  }, []);

  console.log(serverData, isConnected);

  return (
    <div>
      <h1 className="bg-purple-600">Server Statistics</h1>
      {isConnected ? (
        <div>
          {Object.keys(serverData).length > 0 ? (
            <ul>
              {Object.entries(serverData).map(([server, connections]) => (
                <li key={server}>
                  {server}: {connections as any} active connections
                </li>
              ))}
            </ul>
          ) : (
            <p>No server data available</p>
          )}
        </div>
      ) : (
        <p>Connecting to the server...</p>
      )}
    </div>
  );
};

export default ServerStats;
