import { useState, useEffect, createContext, useContext } from "react";
import "./App.css";

// Product descrive la forma dei dati che arrivano dal backend Go.
// Deve combaciare con il JSON restituito da GET /products.
interface Product {
  id: number;
  name: string;
  price: number;
}

// ProductList riceve ancora i prodotti come props: sono dati caricati da App.
// Le azioni del carrello, invece, arrivano dal Context.
interface ProductListProps {
  products: Product[];
}

// CartItem rappresenta una riga del carrello:
// non solo il prodotto, ma anche quante unità di quel prodotto sono nel carrello.
interface CartItem {
  product: Product;
  quantity: number;
}

// CartContextValue è il "contratto" del Context:
// tutti i componenti che leggono il carrello dal Context vedono questi dati
// e possono chiamare queste funzioni.
interface CartContextValue {
  cartItems: CartItem[];
  totalQuantity: number;
  totalPrice: number;
  addToCart: (product: Product) => void;
  decreaseQuantity: (productId: number) => void;
  removeFromCart: (productId: number) => void;
}

// createContext crea un canale condiviso per dati React.
// Usiamo null come valore iniziale perché il valore reale arriverà dal Provider.
const CartContext = createContext<CartContextValue | null>(null);

// useCart è un custom hook: incapsula l'accesso al Context.
// Così ProductList e CartView non devono ripetere useContext + controllo null.
function useCart() {
  const cart = useContext(CartContext);

  if (cart === null) {
    throw new Error("useCart must be used inside CartContext.Provider");
  }

  return cart;
}

function ProductList({ products }: ProductListProps) {
  const { addToCart } = useCart();

  return (
    <ul className="list-group mt-3">
      {/* map trasforma l'array di prodotti in una lista di elementi JSX. */}
      {products.map((product) => (
        <li
          key={product.id}
          className="list-group-item d-flex justify-content-between align-items-center"
        >
          <span>{product.name}</span>

          <div className="d-flex gap-2 align-items-center">
            <strong>€ {product.price}</strong>
            <button
              type="button"
              className="btn btn-sm btn-primary"
              onClick={() => addToCart(product)}
            >
              Aggiungi
            </button>
          </div>
        </li>
      ))}
    </ul>
  );
}

function CartView() {
  const { cartItems, totalPrice, addToCart, decreaseQuantity, removeFromCart } =
    useCart();

  if (cartItems.length === 0) {
    return null;
  }

  return (
    <>
      <h2 className="mt-4">Carrello</h2>

      <ul className="list-group mt-3">
        {cartItems.map((item) => (
          <li
            key={item.product.id}
            className="list-group-item d-flex justify-content-between"
          >
            <span>{item.product.name}</span>
            <div className="d-flex gap-2 align-items-center">
              <button
                type="button"
                className="btn btn-sm btn-outline-secondary"
                onClick={() => decreaseQuantity(item.product.id)}
              >
                -
              </button>

              <strong>Quantità: {item.quantity}</strong>

              <button
                type="button"
                className="btn btn-sm btn-outline-secondary"
                onClick={() => addToCart(item.product)}
              >
                +
              </button>

              <button
                type="button"
                className="btn btn-sm btn-outline-danger"
                onClick={() => removeFromCart(item.product.id)}
              >
                Rimuovi
              </button>
            </div>
          </li>
        ))}
      </ul>

      <p className="status-text">Totale: € {totalPrice.toFixed(2)}</p>
    </>
  );
}

function App() {
  // products è stato locale di App: arriva dal backend e serve alla lista.
  const [products, setProducts] = useState<Product[]>([]);

  // cartItems è lo stato del carrello.
  // Lo terremo nel Context perché più componenti devono leggerlo/modificarlo.
  const [cartItems, setCartItems] = useState<CartItem[]>([]);

  // useEffect con [] parte una volta al montaggio del componente.
  // Qui carichiamo i prodotti dal backend.
  useEffect(() => {
    fetch("http://localhost:8080/products")
      .then((response) => response.json())
      .then((data: Product[]) => setProducts(data))
      .catch((error) => console.error("Errore caricamento prodotti:", error));
  }, []);

  function addToCart(product: Product) {
    // find cerca se il prodotto è già presente nel carrello.
    const existingItem = cartItems.find(
      (item) => item.product.id === product.id,
    );

    if (existingItem) {
      // map crea un nuovo array: aggiorniamo solo la riga corrispondente,
      // lasciando invariate tutte le altre.
      setCartItems(
        cartItems.map((item) =>
          item.product.id === product.id
            ? { ...item, quantity: item.quantity + 1 }
            : item,
        ),
      );
      return;
    }

    // Se il prodotto non era presente, aggiungiamo una nuova riga al carrello.
    setCartItems([...cartItems, { product, quantity: 1 }]);
  }

  // Derived state: non serve uno useState separato.
  // La quantità totale si calcola sempre a partire da cartItems.
  const totalQuantity = cartItems.reduce((sum, item) => {
    return sum + item.quantity;
  }, 0);

  // Altro derived state: prezzo totale = somma di prezzo * quantità.
  const totalPrice = cartItems.reduce((sum, item) => {
    return sum + item.product.price * item.quantity;
  }, 0);

  function decreaseQuantity(productId: number) {
    const existingItem = cartItems.find(
      (item) => item.product.id === productId,
    );

    if (!existingItem) {
      return;
    }

    if (existingItem.quantity === 1) {
      // filter crea un nuovo array senza la riga da rimuovere.
      setCartItems(cartItems.filter((item) => item.product.id !== productId));
      return;
    }

    // Se la quantità è maggiore di 1, decrementiamo solo quella riga.
    setCartItems(
      cartItems.map((item) =>
        item.product.id === productId
          ? { ...item, quantity: item.quantity - 1 }
          : item,
      ),
    );
  }

  function removeFromCart(productId: number) {
    // Rimuove tutta la riga, indipendentemente dalla quantità.
    setCartItems(cartItems.filter((item) => item.product.id !== productId));
  }

  // Oggetto passato al Provider: è il valore che tutti i componenti figli
  // possono leggere tramite useCart().
  const cartContextValue: CartContextValue = {
    cartItems,
    totalQuantity,
    totalPrice,
    addToCart,
    decreaseQuantity,
    removeFromCart,
  };

  return (
    <CartContext.Provider value={cartContextValue}>
      <main className="app-shell">
        <section className="app-panel">
          <p className="eyebrow">04 Global State App</p>
          <h1>Mini carrello</h1>
          <p className="status-text">Prodotti caricati: {products.length}</p>
          <ProductList products={products} />
          <p className="status-text">Articoli nel carrello: {totalQuantity}</p>
          <CartView />
        </section>
      </main>
    </CartContext.Provider>
  );
}

export default App;
