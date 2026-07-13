import { useEffect, useState } from "react";
import "./App.css";

interface Order {
  id: number;
  customer: string;
  status: string;
  total: number;
}
interface OrdersTableProps {
  orders: Order[];
}
interface OrdersResponse {
  items: Order[];
  total: number;
  page: number;
  pageSize: number;
}

interface OrdersFiltersProps {
  search: string;
  status: string;
  onSearchChange: (value: string) => void;
  onStatusChange: (value: string) => void;
}

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

function OrdersFilters({
  search,
  status,
  onSearchChange,
  onStatusChange,
}: OrdersFiltersProps) {
  return (
    <>
      <input
        type="search"
        className="form-control mt-3"
        placeholder="Cerca cliente"
        value={search}
        onChange={(e) => onSearchChange(e.target.value)}
      />

      <select
        className="form-select mt-2"
        value={status}
        onChange={(e) => onStatusChange(e.target.value)}
      >
        <option value="">Tutti gli stati</option>
        <option value="pending">Pending</option>
        <option value="shipped">Shipped</option>
        <option value="delivered">Delivered</option>
      </select>
    </>
  );
}

function App() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(5);
  const totalPages = Math.max(1, Math.ceil(total / pageSize));
  const [search, setSearch] = useState("");
  const [status, setStatus] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [debouncedSearch, setDebouncedSearch] = useState("");

  useEffect(() => {
    const params = new URLSearchParams({
      page: String(page),
      pageSize: String(pageSize),
      search: debouncedSearch,
      status,
    });
    setLoading(true);
    setError("");
    fetch(`http://localhost:8080/orders?${params.toString()}`)
      .then((response) => {
        if (!response.ok) {
          throw new Error("Errore HTTP");
        }

        return response.json();
      })
      .then((data: OrdersResponse) => {
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

  useEffect(() => {
    const timeoutId = window.setTimeout(() => {
      setDebouncedSearch(search);
    }, 400);

    return () => {
      window.clearTimeout(timeoutId);
    };
  }, [search]);

  function handleSearchChange(value: string) {
    setSearch(value);
    setPage(1);
  }

  function handleStatusChange(value: string) {
    setStatus(value);
    setPage(1);
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

        <div className="d-flex gap-2 mt-3">
          <button
            type="button"
            className="btn btn-outline-secondary"
            onClick={() => setPage(page - 1)}
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
            onClick={() => setPage(page + 1)}
            disabled={page >= totalPages}
          >
            Successiva
          </button>
        </div>
      </section>
    </main>
  );
}

export default App;
