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

  useEffect(() => {
    fetch("http://localhost:8080/orders")
      .then((response) => response.json())
      .then((data: OrdersResponse) => {
        setOrders(data.items);
        setTotal(data.total);
      })
      .catch((error) => console.error("Errore caricamento ordini:", error));
  }, []);
  return (
    <main className="app-shell">
      <section className="app-panel">
        <p className="eyebrow">05 Filtered API Table</p>
        <h1>Ordini</h1>
        <p className="status-text">Ordini caricati: {orders.length}</p>
        <p className="status-text">Totale risultati: {total}</p>
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
      </section>
    </main>
  );
}

export default App;
