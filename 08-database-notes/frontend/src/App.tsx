import { useEffect, useState } from "react";
import "./App.css";
import { createNote, deleteNote, fetchNotes, updateNote } from "./api";
import type { Note, NoteRequest } from "./types";
import { NoteForm } from "./components/NoteForm";
import { NotesList } from "./components/NotesList";

function App() {
  const [notes, setNotes] = useState<Note[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [editingNote, setEditingNote] = useState<Note | null>(null);

  // Un solo submit gestisce due casi:
  // - se editingNote è null, creiamo una nuova nota;
  // - se editingNote contiene una nota, aggiorniamo quella nota.
  async function handleSaveNote(payload: NoteRequest) {
    setError("");

    try {
      if (editingNote === null) {
        const createdNote = await createNote(payload);
        setNotes((currentNotes) => [createdNote, ...currentNotes]);
        return;
      }

      const updatedNote = await updateNote(editingNote.id, payload);
      setNotes((currentNotes) =>
        currentNotes.map((note) =>
          note.id === updatedNote.id ? updatedNote : note,
        ),
      );
      setEditingNote(null);
    } catch (error) {
      console.error("Errore salvataggio nota:", error);
      setError("Errore salvataggio nota");
    }
  }

  async function handleDeleteNote(id: number) {
    setError("");

    try {
      await deleteNote(id);
      setNotes((currentNotes) =>
        currentNotes.filter((note) => note.id !== id),
      );

      // Se sto modificando proprio la nota eliminata, esco dalla modalità edit.
      if (editingNote?.id === id) {
        setEditingNote(null);
      }
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

        <NoteForm
          initialValues={
            editingNote === null
              ? undefined
              : {
                  title: editingNote.title,
                  content: editingNote.content,
                }
          }
          submitLabel={editingNote === null ? "Crea nota" : "Salva modifica"}
          onCancel={editingNote === null ? undefined : () => setEditingNote(null)}
          onSubmit={handleSaveNote}
        />

        <NotesList
          notes={notes}
          onEdit={setEditingNote}
          onDelete={handleDeleteNote}
        />
      </section>
    </main>
  );
}

export default App;
