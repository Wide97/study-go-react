import type { OrdersFiltersProps } from "../types";

function OrdersFilters({
  search,
  status,
  onSearchChange,
  onStatusChange,
}: OrdersFiltersProps) {
  return (
    <>
      {/* Input controllato: value arriva dallo state, onChange aggiorna App. */}
      <input
        type="search"
        className="form-control mt-3"
        placeholder="Cerca cliente"
        value={search}
        onChange={(e) => onSearchChange(e.target.value)}
      />

      {/* Anche la select è controllata dallo state status. */}
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

export { OrdersFilters };
