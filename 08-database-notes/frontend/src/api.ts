import type { Note, NoteRequest } from "./types";

const API_BASE_URL = "http://localhost:8080";

// api.ts è lo strato che conosce gli URL del backend.
// I componenti React chiamano queste funzioni invece di scrivere fetch ovunque.
export async function fetchNotes(): Promise<Note[]> {
  const response = await fetch(`${API_BASE_URL}/notes`);

  if (!response.ok) {
    throw new Error("Failed to fetch notes");
  }

  return response.json();
}

export async function createNote(payload: NoteRequest): Promise<Note> {
  const response = await fetch(`${API_BASE_URL}/notes`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    // Il body HTTP può essere testo o byte.
    // JSON.stringify trasforma l'oggetto TypeScript in una stringa JSON.
    body: JSON.stringify(payload),
  });

  if (!response.ok) {
    throw new Error("Failed to create note");
  }

  return response.json();
}

export async function deleteNote(id: number): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/notes/${id}`, {
    method: "DELETE",
  });

  if (!response.ok) {
    throw new Error("Failed to delete note");
  }

  // Il backend risponde 204 No Content.
  // Non chiamiamo response.json() perché il body è vuoto.
}

export async function updateNote(
  id: number,
  payload: NoteRequest,
): Promise<Note> {
  const response = await fetch(`${API_BASE_URL}/notes/${id}`, {
    method: "PUT",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });

  if (!response.ok) {
    throw new Error("Failed to update note");
  }

  return response.json();
}
