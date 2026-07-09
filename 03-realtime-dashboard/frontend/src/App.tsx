import { useEffect, useState } from "react";
import "./App.css";

interface SensorMessage {
  value: number;
}

function App() {
  const [value, setValue] = useState<number | null>(null);
  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8080/ws");

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
        <p className="status-text">
          {value === null ? "In attesa di dati..." : `Valore live: ${value}`}
        </p>
      </section>
    </main>
  );
}

export default App;
