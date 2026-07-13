import { useEffect, useState } from "react";
import "./App.css";
import { fetchOrders } from "./api";
import type { Order, OrdersTableProps, PaginationControlsProps } from "./types";
import { OrdersFilters } from "./components/OrdersFilters";

function OrdersTable({ orders }: OrdersTableProps) {
  return (
    <table className="table table-sm mt-3">
      <thead>
        <tr>
          <th>ID</th>
          <th>Cliente</th>
          <th>Stato</th>
          <th>Totale</th>
        </tr>
      </thead>

      <tbody>
        {orders.length === 0 && (
          <tr>
            <td colSpan={4} className="text-center text-muted">
              Nessun ordine trovato
            </td>
          </tr>
        )}

        {/* map trasforma l'array di ordini in righe di tabella JSX. */}
        {orders.map((order) => (
          <tr key={order.id}>
            <td>{order.id}</td>
            <td>{order.customer}</td>
            <td>{order.status}</td>
            <td>€ {order.total.toFixed(2)}</td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

function PaginationControls({
  page,
  totalPages,
  onPrevious,
  onNext,
}: PaginationControlsProps) {
  return (
    <div className="d-flex gap-2 mt-3">
      <button
        type="button"
        className="btn btn-outline-secondary"
        onClick={onPrevious}
        // Prima pagina: non ha senso andare indietro.
        disabled={page === 1}
      >
        Precedente
      </button>

      <span className="align-self-center">
        Pagina {page} di {totalPages}
      </span>

      <button
        type="button"
        className="btn btn-outline-secondary"
        onClick={onNext}
        // Ultima pagina: non ha senso andare avanti.
        disabled={page >= totalPages}
      >
        Successiva
      </button>
    </div>
  );
}

function App() {
  // orders contiene solo la pagina corrente ricevuta dal backend.
  const [orders, setOrders] = useState<Order[]>([]);

  // total è il totale dei risultati dopo filtri, prima della paginazione.
  // Serve per calcolare quante pagine esistono.
  const [total, setTotal] = useState(0);

  // page è la pagina richiesta al backend.
  const [page, setPage] = useState(1);

  // Per ora pageSize è fisso: lo teniamo in state perché fa parte della query.
  const [pageSize] = useState(5);

  // Derived state: non serve salvarlo con useState perché deriva da total/pageSize.
  const totalPages = Math.max(1, Math.ceil(total / pageSize));

  // search è il valore immediato scritto nell'input.
  const [search, setSearch] = useState("");

  // status è il filtro stato selezionato.
  const [status, setStatus] = useState("");

  // loading/error rappresentano gli stati della richiesta HTTP.
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  // debouncedSearch è il valore usato davvero nella fetch.
  // Cambia solo dopo che l'utente smette di scrivere per 400ms.
  const [debouncedSearch, setDebouncedSearch] = useState("");

  // Questa useEffect carica gli ordini ogni volta che cambiano pagina,
  // dimensione pagina, ricerca debounced o filtro stato.
  useEffect(() => {
    setLoading(true);
    setError("");

    fetchOrders({
      page,
      pageSize,
      search: debouncedSearch,
      status,
    })
      .then((data) => {
        setOrders(data.items);
        setTotal(data.total);
        setLoading(false);
      })
      .catch((error) => {
        console.error("Errore caricamento ordini:", error);
        setError("Errore caricamento ordini");
        setLoading(false);
      });
  }, [page, pageSize, debouncedSearch, status]);

  // Questa useEffect implementa il debounce della ricerca.
  // Ogni volta che search cambia, parte un timer. Se l'utente scrive ancora,
  // il cleanup cancella il timer precedente.
  useEffect(() => {
    const timeoutId = window.setTimeout(() => {
      setDebouncedSearch(search);
    }, 400);

    return () => {
      window.clearTimeout(timeoutId);
    };
  }, [search]);

  function handleSearchChange(value: string) {
    // Cambiando filtro torniamo alla prima pagina: la pagina corrente potrebbe
    // non esistere più con i nuovi risultati.
    setSearch(value);
    setPage(1);
  }

  function handleStatusChange(value: string) {
    // Stessa logica della ricerca: nuovo filtro, si riparte da pagina 1.
    setStatus(value);
    setPage(1);
  }

  function goToPreviousPage() {
    setPage(page - 1);
  }

  function goToNextPage() {
    setPage(page + 1);
  }

  return (
    <main className="app-shell">
      <section className="app-panel">
        <p className="eyebrow">05 Filtered API Table</p>
        <h1>Ordini</h1>
        <p className="status-text">Ordini caricati: {orders.length}</p>
        <p className="status-text">Totale risultati: {total}</p>
        <OrdersFilters
          search={search}
          status={status}
          onSearchChange={handleSearchChange}
          onStatusChange={handleStatusChange}
        />
        {loading && <div className="alert alert-info mt-3">Caricamento...</div>}

        {error !== "" && <div className="alert alert-danger mt-3">{error}</div>}

        <OrdersTable orders={orders} />

        <PaginationControls
          page={page}
          totalPages={totalPages}
          onPrevious={goToPreviousPage}
          onNext={goToNextPage}
        />
      </section>
    </main>
  );
}

export default App;
