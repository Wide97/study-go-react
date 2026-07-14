import type { Note } from "../types";

interface NotesListProps {
  notes: Note[];
  onEdit: (note: Note) => void;
  onDelete: (id: number) => void;
}

export function NotesList({ notes, onEdit, onDelete }: NotesListProps) {
  return (
    <ul className="list-group mt-3">
      {notes.map((note) => (
        <li key={note.id} className="list-group-item">
          <div className="d-flex justify-content-between gap-2">
            <div>
              <h2 className="h5 mb-1">{note.title}</h2>
              <p className="mb-0">{note.content}</p>
            </div>

            <div className="d-flex gap-2 align-items-start">
              <button
                type="button"
                className="btn btn-sm btn-outline-secondary"
                onClick={() => onEdit(note)}
              >
                Modifica
              </button>

              <button
                type="button"
                className="btn btn-sm btn-outline-danger"
                onClick={() => onDelete(note.id)}
              >
                Elimina
              </button>
            </div>
          </div>
        </li>
      ))}
    </ul>
  );
}
