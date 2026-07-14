import type { Note } from "../types";

interface NotesListProps {
  notes: Note[];
}

export function NotesList({ notes }: NotesListProps) {
  return (
    <ul className="list-group mt-3">
      {notes.map((note) => (
        <li key={note.id} className="list-group-item">
          <h2 className="h5 mb-1">{note.title}</h2>
          <p className="mb-0">{note.content}</p>
        </li>
      ))}
    </ul>
  );
}
