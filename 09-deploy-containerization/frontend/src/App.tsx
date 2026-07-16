import { useEffect, useState } from "react";
import "./App.css";

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080";

function App() {
  const [status, setStatus] = useState("checking");
  const [message, setMessage] = useState("");

  useEffect(() => {
    // Il frontend non conosce Docker direttamente: conosce solo un URL API.
    // In locale puo essere http://localhost:8080.
    // In container useremo nginx per inoltrare /api verso il backend.
    fetch(`${API_BASE_URL}/health`)
      .then((response) => {
        if (!response.ok) {
          throw new Error("Backend not healthy");
        }

        return response.text();
      })
      .then((text) => {
        setStatus("connected");
        setMessage(text);
      })
      .catch((error) => {
        console.error("Backend health check failed:", error);
        setStatus("error");
        setMessage("Backend non raggiungibile");
      });
  }, []);

  return (
    <main className="app-shell">
      <section className="app-panel">
        <p className="eyebrow">09 Deploy Containerization</p>
        <h1>Full-stack Docker</h1>
        <p className="status-text">Backend status: {status}</p>
        <p className="status-text">Risposta: {message || "..."}</p>
      </section>
    </main>
  );
}

export default App;
