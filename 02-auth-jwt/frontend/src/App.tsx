import { useEffect, useState } from "react";
import "./App.css";

// Base URL del backend Go.
// Tenerla in una costante evita di ripetere la stringa in più fetch.
const API_URL = "http://localhost:8080";

// Chiave usata per salvare il token nel localStorage del browser.
// localStorage conserva il valore anche dopo refresh della pagina.
const TOKEN_STORAGE_KEY = "authToken";

// Forma del JSON restituito da POST /login.
// Il backend risponde con { "token": "..." }.
interface LoginResponse {
  token: string;
}

// Forma del JSON restituito da GET /me.
// Il backend legge il JWT e risponde con l'email dentro il claim sub.
interface MeResponse {
  email: string;
}

function App() {
  // Stati dei campi del form.
  // Sono input controllati: il valore visualizzato viene da React,
  // e ogni onChange aggiorna lo stato corrispondente.
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  // Stato del token JWT.
  // L'inizializzazione usa una funzione invece di leggere localStorage
  // direttamente come valore: React la esegue solo al primo render.
  const [token, setToken] = useState(() => localStorage.getItem(TOKEN_STORAGE_KEY) ?? "");

  // Email dell'utente autenticato, caricata da GET /me.
  const [userEmail, setUserEmail] = useState("");

  // Messaggio di errore mostrato nella UI.
  const [error, setError] = useState("");

  // Stati booleani per disabilitare/mostrare feedback durante chiamate async.
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isLoadingProfile, setIsLoadingProfile] = useState(false);

  // Derivato dallo stato token: se c'è un token, consideriamo l'utente loggato.
  // Non è uno useState separato perché si può calcolare sempre da token.
  const isAuthenticated = token !== "";

  // useEffect esegue codice "a lato" del render.
  // Qui serve per caricare automaticamente il profilo quando il token cambia:
  // - dopo login;
  // - oppure al refresh, se un token era già salvato nel localStorage.
  useEffect(() => {
    if (token === "") {
      return;
    }

    // fetchMe è async, ma la funzione passata a useEffect non deve essere async.
    // Usiamo void per indicare esplicitamente che non aspettiamo la Promise qui.
    void fetchMe(token);
  }, [token]);

  // Gestisce il submit del form di login.
  // Il form HTML ricaricherebbe la pagina di default: preventDefault lo blocca,
  // così possiamo gestire tutto via React/fetch.
  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();

    // Prima di una nuova login puliamo eventuali messaggi/stati vecchi.
    setError("");
    setUserEmail("");
    setIsSubmitting(true);

    try {
      // fetch invia la richiesta HTTP al backend.
      // Il body deve essere una stringa JSON, quindi usiamo JSON.stringify.
      const response = await fetch(`${API_URL}/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });

      // response.ok è true solo per status 2xx.
      // Il backend risponde 401 se email/password non sono valide.
      if (!response.ok) {
        setError("Email o password non valide");
        return;
      }

      // response.json() legge il body JSON e lo converte in oggetto JS.
      const data: LoginResponse = await response.json();

      // Salviamo il token sia in localStorage sia nello stato React:
      // - localStorage mantiene la sessione dopo refresh;
      // - setToken fa ridisegnare subito la UI.
      localStorage.setItem(TOKEN_STORAGE_KEY, data.token);
      setToken(data.token);

      // Puliamo la password dal form dopo un login riuscito.
      setPassword("");
    } catch {
      // Qui finiamo per errori di rete, backend spento, CORS non raggiungibile,
      // o problemi simili prima di ricevere una risposta HTTP valida.
      setError("Backend non raggiungibile");
    } finally {
      // finally gira sia in caso di successo sia in caso di errore.
      setIsSubmitting(false);
    }
  }

  // Chiama la rotta protetta GET /me.
  // Riceve il token come parametro, invece di leggere direttamente lo stato,
  // così può essere usata anche subito dopo il refresh da useEffect.
  async function fetchMe(authToken: string) {
    setError("");
    setIsLoadingProfile(true);

    try {
      // Il token JWT va mandato nell'header Authorization con schema Bearer.
      const response = await fetch(`${API_URL}/me`, {
        headers: {
          Authorization: `Bearer ${authToken}`,
        },
      });

      // Se /me risponde 401, il token è assente, invalido o scaduto.
      // In quel caso puliamo localStorage e stato: l'utente torna al login.
      if (!response.ok) {
        localStorage.removeItem(TOKEN_STORAGE_KEY);
        setToken("");
        setUserEmail("");
        setError("Sessione non valida o scaduta");
        return;
      }

      const data: MeResponse = await response.json();
      setUserEmail(data.email);
    } catch {
      setError("Backend non raggiungibile");
    } finally {
      setIsLoadingProfile(false);
    }
  }

  // Logout lato frontend.
  // Con JWT stateless il backend non ha una sessione da distruggere:
  // basta eliminare il token dal client.
  function handleLogout() {
    localStorage.removeItem(TOKEN_STORAGE_KEY);
    setToken("");
    setUserEmail("");
    setPassword("");
    setError("");
  }

  return (
    <main className="app-shell">
      <section className="app-panel">
        <h1>{isAuthenticated ? "Area protetta" : "Login"}</h1>

        {/* Rendering condizionale: se non siamo autenticati mostriamo il form,
            altrimenti mostriamo la schermata protetta. */}
        {!isAuthenticated ? (
          <form onSubmit={handleSubmit} className="mt-4">
            <div className="mb-3">
              <label htmlFor="email" className="form-label">
                Email
              </label>
              <input
                id="email"
                type="email"
                className="form-control"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />
            </div>

            <div className="mb-3">
              <label htmlFor="password" className="form-label">
                Password
              </label>
              <input
                id="password"
                type="password"
                className="form-control"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
              />
            </div>

            {/* disabled evita doppi submit mentre la richiesta è in corso. */}
            <button type="submit" className="btn btn-primary" disabled={isSubmitting}>
              {isSubmitting ? "Accesso..." : "Login"}
            </button>
          </form>
        ) : (
          <div className="protected-view">
            <div className="alert alert-success mb-0">Login riuscito</div>

            <div className="profile-box">
              <span>Utente autenticato</span>
              {/* Durante la chiamata /me mostriamo un testo temporaneo.
                  Poi mostriamo l'email restituita dal backend. */}
              <strong>{isLoadingProfile ? "Caricamento..." : userEmail}</strong>
            </div>

            <button type="button" className="btn btn-outline-secondary" onClick={handleLogout}>
              Logout
            </button>
          </div>
        )}

        {/* Mostriamo l'errore solo se la stringa non è vuota. */}
        {error !== "" && (
          <div className="alert alert-danger mt-3 mb-0">
            {error}
          </div>
        )}
      </section>
    </main>
  );
}

export default App;
