import {
  cleanup,
  fireEvent,
  render,
  screen,
  waitFor,
} from "@testing-library/react";
import { afterEach, describe, expect, test, vi } from "vitest";
import App from "./App";

interface MockOrder {
  id: number;
  customer: string;
  status: string;
  total: number;
}

interface MockOrdersResponse {
  items: MockOrder[];
  total: number;
  page: number;
  pageSize: number;
}

function createOrdersResponse(items: MockOrder[]): MockOrdersResponse {
  return {
    items,
    total: items.length,
    page: 1,
    pageSize: 5,
  };
}

function mockFetchSuccess(response: MockOrdersResponse) {
  // In questi test non vogliamo chiamare davvero http://localhost:8080.
  // Sostituiamo globalThis.fetch con una funzione finta controllata dal test.
  const fetchMock = vi.fn().mockResolvedValue({
    ok: true,
    json: async () => response,
  });

  vi.stubGlobal("fetch", fetchMock);

  return fetchMock;
}

function mockFetchError() {
  const fetchMock = vi.fn().mockResolvedValue({
    ok: false,
    json: async () => ({}),
  });

  vi.stubGlobal("fetch", fetchMock);

  return fetchMock;
}

afterEach(() => {
  // Ogni test deve lasciare l'ambiente pulito.
  // Se non ripristiniamo i mock, un test può influenzare quello successivo.
  cleanup();
  vi.restoreAllMocks();
  vi.unstubAllGlobals();
  vi.useRealTimers();
});

describe("App", () => {
  test("carica e mostra gli ordini restituiti dal backend", async () => {
    mockFetchSuccess(
      createOrdersResponse([
        { id: 1, customer: "Alice", status: "pending", total: 150 },
        { id: 2, customer: "Bob", status: "shipped", total: 200 },
      ]),
    );

    render(<App />);

    // findByText aspetta che l'elemento compaia.
    // È utile quando la UI dipende da una fetch asincrona.
    expect(await screen.findByText("Alice")).toBeInTheDocument();
    expect(screen.getByText("Bob")).toBeInTheDocument();
    expect(screen.getByText("Totale risultati: 2")).toBeInTheDocument();
  });

  test("mostra un messaggio quando non ci sono ordini", async () => {
    mockFetchSuccess(createOrdersResponse([]));

    render(<App />);

    expect(await screen.findByText("Nessun ordine trovato")).toBeInTheDocument();
  });

  test("aggiunge il filtro status alla query string", async () => {
    const fetchMock = mockFetchSuccess(createOrdersResponse([]));

    render(<App />);
    await screen.findByText("Nessun ordine trovato");

    fireEvent.change(screen.getByRole("combobox"), {
      target: { value: "pending" },
    });

    await waitFor(() => {
      expect(fetchMock).toHaveBeenLastCalledWith(
        "http://localhost:8080/orders?page=1&pageSize=5&search=&status=pending",
      );
    });
  });

  test("applica il debounce prima di cercare", async () => {
    const fetchMock = mockFetchSuccess(createOrdersResponse([]));

    render(<App />);
    await screen.findByText("Nessun ordine trovato");

    fireEvent.change(screen.getByPlaceholderText("Cerca cliente"), {
      target: { value: "Ali" },
    });

    // Subito dopo la digitazione non parte una nuova fetch:
    // il debounce deve aspettare 400ms.
    expect(fetchMock).toHaveBeenCalledTimes(1);

    // Usiamo un timer reale per rendere il test più semplice da leggere.
    // Dopo 450ms il debounce da 400ms deve aver aggiornato debouncedSearch.
    await new Promise((resolve) => window.setTimeout(resolve, 450));

    await waitFor(() => {
      expect(fetchMock).toHaveBeenCalledTimes(2);
    });

    expect(fetchMock).toHaveBeenLastCalledWith(
      "http://localhost:8080/orders?page=1&pageSize=5&search=Ali&status=",
    );
  });

  test("mostra un errore quando la risposta HTTP non è ok", async () => {
    mockFetchError();

    render(<App />);

    expect(await screen.findByText("Errore caricamento ordini")).toBeInTheDocument();
  });
});
