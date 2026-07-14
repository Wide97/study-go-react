  export interface Note {
    id: number;
    title: string;
    content: string;
    created_at: string;
    updated_at: string;
  }

  export interface NoteRequest {
    title: string;
    content: string;
  }