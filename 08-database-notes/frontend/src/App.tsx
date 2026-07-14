import { useEffect, useState } from "react";
import "./App.css";
import { fetchNotes, createNote, deleteNote } from "./api";
import type { Note, NoteRequest } from "./types";
import { NotesList } from "./components/NotesList";
import { NoteForm } from "./components/NoteForm";

function App() {
  const [notes, setNotes] = useState<Note[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  async function handleCreateNote(payload: NoteRequest) {
    setError("");

    try {
      const createdNote = await createNote(payload);
      setNotes([createdNote, ...notes]);
    } catch (error) {
      console.error("Errore creazione nota:", error);
      setError("Errore creazione nota");
    }
  }

  async function handleDeleteNote(id: number) {
    setError("");

    try {
      await deleteNote(id);
      setNotes(notes.filter((note) => note.id !== id));
    } catch (error) {
      console.error("Errore eliminazione nota:", error);
      setError("Errore eliminazione nota");
    }
  }

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
        <NoteForm onSubmit={handleCreateNote} />
        <NotesList notes={notes} onDelete={handleDeleteNote} />
      </section>
    </main>
  );
}

export default App;
