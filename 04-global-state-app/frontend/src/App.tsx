import { useState, useEffect, createContext, useContext } from "react";
import "./App.css";

interface Product {
  id: number;
  name: string;
  price: number;
}

interface ProductListProps {
  products: Product[];
}

interface CartItem {
  product: Product;
  quantity: number;
}

interface CartContextValue {
  cartItems: CartItem[];
  totalQuantity: number;
  totalPrice: number;
  addToCart: (product: Product) => void;
  decreaseQuantity: (productId: number) => void;
  removeFromCart: (productId: number) => void;
}

const CartContext = createContext<CartContextValue | null>(null);

function ProductList({ products }: ProductListProps) {
  const cart = useContext(CartContext);

  if (cart === null) {
    return null;
  }

  const { addToCart } = cart;

  return (
    <ul className="list-group mt-3">
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
  const cart = useContext(CartContext);

  if (cart === null) {
    return null;
  }

  const { cartItems, totalPrice, addToCart, decreaseQuantity, removeFromCart } =
    cart;

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
  const [products, setProducts] = useState<Product[]>([]);
  const [cartItems, setCartItems] = useState<CartItem[]>([]);
  useEffect(() => {
    fetch("http://localhost:8080/products")
      .then((response) => response.json())
      .then((data: Product[]) => setProducts(data))
      .catch((error) => console.error("Errore caricamento prodotti:", error));
  }, []);

  function addToCart(product: Product) {
    const existingItem = cartItems.find(
      (item) => item.product.id === product.id,
    );

    if (existingItem) {
      setCartItems(
        cartItems.map((item) =>
          item.product.id === product.id
            ? { ...item, quantity: item.quantity + 1 }
            : item,
        ),
      );
      return;
    }

    setCartItems([...cartItems, { product, quantity: 1 }]);
  }

  const totalQuantity = cartItems.reduce((sum, item) => {
    return sum + item.quantity;
  }, 0);

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
      setCartItems(cartItems.filter((item) => item.product.id !== productId));
      return;
    }

    setCartItems(
      cartItems.map((item) =>
        item.product.id === productId
          ? { ...item, quantity: item.quantity - 1 }
          : item,
      ),
    );
  }

  function removeFromCart(productId: number) {
    setCartItems(cartItems.filter((item) => item.product.id !== productId));
  }

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
