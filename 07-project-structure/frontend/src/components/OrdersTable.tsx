import type { OrdersTableProps } from "../types";

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

export { OrdersTable };
