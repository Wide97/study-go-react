import { useEffect, useState } from "react";
import "./App.css";

interface Order {
  id: number;
  customer: string;
  status: string;
  total: number;
}

interface OrdersResponse {
  items: Order[];
  total: number;
  page: number;
  pageSize: number;
}

function App() {
  const [orders, setOrders] = useState<Order[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [pageSize] = useState(5);
  const totalPages = Math.max(1, Math.ceil(total / pageSize));
  const [search, setSearch] = useState("");
  const [status, setStatus] = useState("");

  useEffect(() => {
    fetch(
      `http://localhost:8080/orders?page=${page}&pageSize=${pageSize}&search=${search}&status=${status}`,
    )
      .then((response) => response.json())
      .then((data: OrdersResponse) => {
        setOrders(data.items);
        setTotal(data.total);
      })
      .catch((error) => console.error("Errore caricamento ordini:", error));
  }, [page, pageSize, search, status]);

  return (
    <main className="app-shell">
      <section className="app-panel">
        <p className="eyebrow">05 Filtered API Table</p>
        <h1>Ordini</h1>
        <p className="status-text">Ordini caricati: {orders.length}</p>
        <p className="status-text">Totale risultati: {total}</p>
        <input
          type="search"
          className="form-control mt-3"
          placeholder="Cerca cliente"
          value={search}
          onChange={(e) => {
            setSearch(e.target.value);
            setPage(1);
          }}
        />
        <select
          className="form-select mt-2"
          value={status}
          onChange={(e) => {
            setStatus(e.target.value);
            setPage(1);
          }}
        >
          <option value="">Tutti gli stati</option>
          <option value="pending">Pending</option>
          <option value="shipped">Shipped</option>
          <option value="delivered">Delivered</option>
        </select>
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
