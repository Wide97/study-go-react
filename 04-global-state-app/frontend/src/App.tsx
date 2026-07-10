import { useState, useEffect } from "react";
import "./App.css";

interface Product {
  id: number;
  name: string;
  price: number;
}

function App() {
  const [products, setProducts] = useState<Product[]>([]);

  useEffect(() => {
    fetch("http://localhost:8080/products")
      .then((response) => response.json())
      .then((data: Product[]) => setProducts(data))
      .catch((error) => console.error("Errore caricamento prodotti:", error));
  }, []);

  return (
    <main className="app-shell">
      <section className="app-panel">
        <p className="eyebrow">04 Global State App</p>
        <h1>Mini carrello</h1>
        <p className="status-text">Prodotti caricati: {products.length}</p>{" "}
      </section>
    </main>
  );
}

export default App;
