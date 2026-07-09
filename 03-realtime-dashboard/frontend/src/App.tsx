import { useEffect, useState } from "react";
import "./App.css";

interface SensorMessage {
  value: number;
}

type ConnectionStatus = "connecting" | "connected" | "disconnected" | "error";

function App() {
  const [value, setValue] = useState<number | null>(null);
  const [status, setStatus] = useState<ConnectionStatus>("connecting");
  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8080/ws");

    socket.onopen = () => {
      setStatus("connected");
    };

    socket.onclose = () => {
      setStatus("disconnected");
    };

    socket.onerror = () => {
      setStatus("error");
    };

    socket.onmessage = (event) => {
      const message: SensorMessage = JSON.parse(event.data);
      setValue(message.value);
    };

    return () => {
      socket.close();
    };
  }, []);

  return (
    <main className="app-shell">
      <section className="app-panel">
        <p className="eyebrow">03 Realtime Dashboard</p>
        <h1>Dashboard realtime</h1>
        <p className="status-text">Stato: {status}</p>
        <p className="status-text">
          {value === null ? "In attesa di dati..." : `Valore live: ${value}`}
        </p>
      </section>
    </main>
  );
}

export default App;
