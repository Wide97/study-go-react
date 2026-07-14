import type { Note, NoteRequest } from "./types";

const API_BASE_URL = "http://localhost:8080";

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
      body: JSON.stringify(payload),
    });

    if (!response.ok) {
      throw new Error("Failed to create note");
    }

    return response.json();
  }
