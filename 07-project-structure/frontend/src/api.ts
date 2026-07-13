import type { OrdersResponse } from "./types";

interface FetchOrdersParams {
  page: number;
  pageSize: number;
  search: string;
  status: string;
}

// fetchOrders incapsula la chiamata HTTP al backend.
// App passa parametri applicativi e riceve una risposta tipizzata, senza
// doversi occupare di costruire l'URL a mano.
export async function fetchOrders({
  page,
  pageSize,
  search,
  status,
}: FetchOrdersParams): Promise<OrdersResponse> {
  const params = new URLSearchParams({
    page: String(page),
    pageSize: String(pageSize),
    search,
    status,
  });

  const response = await fetch(
    `http://localhost:8080/orders?${params.toString()}`,
  );

  // fetch entra nel catch solo per errori di rete.
  // Con response.ok gestiamo anche risposte HTTP come 404 o 500.
  if (!response.ok) {
    throw new Error("Errore HTTP");
  }

  return response.json();
}
