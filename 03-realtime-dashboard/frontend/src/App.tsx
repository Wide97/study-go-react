import { useEffect, useState } from "react";
import "./App.css";

// Forma del messaggio JSON mandato dal backend via WebSocket.
// Il backend invia oggetti come {"value":1}, {"value":2}, ...
interface SensorMessage {
  value: number;
}

// Union type: status può essere solo uno di questi quattro valori.
// Così TypeScript ci protegge da typo tipo "connectd" o "errore".
type ConnectionStatus = "connecting" | "connected" | "disconnected" | "error";

function App() {
  // Ultimo valore ricevuto dal WebSocket.
  // null significa: non abbiamo ancora ricevuto nessun messaggio.
  const [value, setValue] = useState<number | null>(null);

  // Stato della connessione WebSocket, mostrato nella UI.
  const [status, setStatus] = useState<ConnectionStatus>("connecting");

  // useEffect serve per aprire la connessione WebSocket dopo il primo render.
  // L'array vuoto [] significa: esegui questo codice una sola volta, quando
  // il componente viene montato.
  useEffect(() => {
    // WebSocket è un'API del browser.
    // Usiamo ws:// invece di http:// perché il protocollo è diverso.
    const socket = new WebSocket("ws://localhost:8080/ws");

    // onopen scatta quando il browser ha completato la connessione col backend.
    socket.onopen = () => {
      setStatus("connected");
    };

    // onclose scatta quando la connessione viene chiusa.
    // Può succedere se il server si ferma, se chiudiamo la pagina,
    // o quando il cleanup dell'effect chiama socket.close().
    socket.onclose = () => {
      setStatus("disconnected");
    };

    // onerror scatta quando il browser rileva un problema sulla connessione.
    socket.onerror = () => {
      setStatus("error");
    };

    // onmessage scatta ogni volta che il backend invia un messaggio.
    // event.data arriva come stringa, quindi la convertiamo con JSON.parse.
    socket.onmessage = (event) => {
      const message: SensorMessage = JSON.parse(event.data);
      setValue(message.value);
    };

    // Cleanup dell'effect: React lo esegue quando il componente viene smontato.
    // Chiudere il WebSocket evita connessioni lasciate aperte.
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
          {/* Rendering condizionale: prima del primo messaggio mostriamo attesa,
              dopo mostriamo il valore live ricevuto dal backend. */}
          {value === null ? "In attesa di dati..." : `Valore live: ${value}`}
        </p>
      </section>
    </main>
  );
}

export default App;
