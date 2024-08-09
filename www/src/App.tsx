import { useEffect, useState } from "react";
import "./App.css";
import BackendList from "./components/BackendList";

function App() {
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
      {/* <ServerStats /> */}
      {isConnected ? (
        <div>
          <BackendList data={serverData} />
        </div>
      ) : (
        <p>Connecting to the server...</p>
      )}
    </div>
  );
}

export default App;
