// Order descrive un singolo ordine ricevuto dal backend.
// Deve combaciare con i campi JSON restituiti da Go: id, customer, status, total.
export interface Order {
  id: number;
  customer: string;
  status: string;
  total: number;
}

// OrdersTable riceve solo gli ordini da mostrare.
// Non sa nulla di fetch, filtri o paginazione: il suo compito è solo renderizzare
// una tabella.
export interface OrdersTableProps {
  orders: Order[];
}

// OrdersResponse descrive la risposta completa di GET /orders.
// items è la pagina corrente, total è il totale dopo i filtri, page/pageSize
// descrivono la paginazione usata dal backend.
export interface OrdersResponse {
  items: Order[];
  total: number;
  page: number;
  pageSize: number;
}

// OrdersFilters riceve valori e callback da App.
// Gli input sono controllati: il valore visibile arriva dallo state React.
export interface OrdersFiltersProps {
  search: string;
  status: string;
  onSearchChange: (value: string) => void;
  onStatusChange: (value: string) => void;
}

// PaginationControls non modifica direttamente lo state.
// Espone solo eventi: App decide cosa fare quando l'utente clicca.
export interface PaginationControlsProps {
  page: number;
  totalPages: number;
  onPrevious: () => void;
  onNext: () => void;
}