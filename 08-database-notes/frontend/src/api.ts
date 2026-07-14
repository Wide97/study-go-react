import type { Note } from "./types";

const API_BASE_URL = "http://localhost:8080";

  export async function fetchNotes(): Promise<Note[]> {
    const response = await fetch(`${API_BASE_URL}/notes`);

    if (!response.ok) {
      throw new Error("Failed to fetch notes");
    }

    return response.json();
  }