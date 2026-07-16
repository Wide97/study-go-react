import { useEffect, useState } from "react";
import "./App.css";

// Il browser chiama sempre lo stesso origin tramite /api.
// In sviluppo Vite inoltra /api al backend locale; nel container nginx lo
// inoltra al servizio backend. In questo modo evitiamo CORS e non inseriamo
// host Docker nel codice React.
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? "/api";

function App() {
  const [status, setStatus] = useState("checking");
  const [message, setMessage] = useState("");

  useEffect(() => {
    // Il frontend non conosce Docker direttamente: conosce solo /api.
    // Il server che serve il frontend sa come raggiungere il backend.
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
