import { useState } from "react";
import "./App.css";

function App() {
  const [value] = useState<number | null>(null);
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
