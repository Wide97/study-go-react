import { useState } from "react";
import type { NoteRequest } from "../types";

interface NoteFormProps {
  onSubmit: (payload: NoteRequest) => void;
}

export function NoteForm({ onSubmit }: NoteFormProps) {
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");

  function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();

    onSubmit({
      title,
      content,
    });

    setTitle("");
    setContent("");
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

      <button type="submit" className="btn btn-primary">
        Crea nota
      </button>
    </form>
  );
}
