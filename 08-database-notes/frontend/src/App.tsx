import { useEffect, useState } from "react";
import "./App.css";
import { fetchNotes } from "./api";
import type { Note } from "./types";
import { NotesList } from "./components/NotesList";

function App() {
  const [notes, setNotes] = useState<Note[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    setLoading(true);
    setError("");

    fetchNotes()
      .then((data) => {
        setNotes(data);
        setLoading(false);
      })
      .catch((error) => {
        console.error("Errore caricamento note:", error);
        setError("Errore caricamento note");
        setLoading(false);
      });
  }, []);

  return (
    <main className="app-shell">
      <section className="app-panel">
        <p className="eyebrow">08 Database Notes</p>
        <h1>Note</h1>
        <p className="status-text">Note caricate: {notes.length}</p>

        {loading && <div className="alert alert-info mt-3">Caricamento...</div>}

        {error !== "" && <div className="alert alert-danger mt-3">{error}</div>}

        <NotesList notes={notes} />
      </section>
    </main>
  );
}

export default App;
