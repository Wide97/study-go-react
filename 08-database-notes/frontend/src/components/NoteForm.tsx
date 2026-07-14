import { useEffect, useState } from "react";
import type { NoteRequest } from "../types";

interface NoteFormProps {
  initialValues?: NoteRequest;
  submitLabel: string;
  onCancel?: () => void;
  onSubmit: (payload: NoteRequest) => void;
}

export function NoteForm({
  initialValues,
  submitLabel,
  onCancel,
  onSubmit,
}: NoteFormProps) {
  const [title, setTitle] = useState(initialValues?.title ?? "");
  const [content, setContent] = useState(initialValues?.content ?? "");

  // useState usa initialValues solo al primo render.
  // Quando l'utente sceglie un'altra nota da modificare, initialValues cambia:
  // questo effect sincronizza i campi del form con la nota selezionata.
  useEffect(() => {
    setTitle(initialValues?.title ?? "");
    setContent(initialValues?.content ?? "");
  }, [initialValues]);

  function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();

    onSubmit({
      title,
      content,
    });

    // Dopo una creazione svuotiamo il form.
    // In modifica, invece, App azzera editingNote e l'effect sopra pulisce i campi.
    if (initialValues === undefined) {
      setTitle("");
      setContent("");
    }
  }

  return (
    <form className="mt-3" onSubmit={handleSubmit}>
      <div className="mb-2">
        <label htmlFor="title" className="form-label">
          Titolo
        </label>
        <input
          id="title"
          className="form-control"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
        />
      </div>

      <div className="mb-2">
        <label htmlFor="content" className="form-label">
          Contenuto
        </label>
        <textarea
          id="content"
          className="form-control"
          value={content}
          onChange={(e) => setContent(e.target.value)}
        />
      </div>

      <div className="d-flex gap-2">
        <button type="submit" className="btn btn-primary">
          {submitLabel}
        </button>

        {onCancel !== undefined && (
          <button
            type="button"
            className="btn btn-outline-secondary"
            onClick={onCancel}
          >
            Annulla
          </button>
        )}
      </div>
    </form>
  );
}
