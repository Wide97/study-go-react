import { useEffect, useState } from "react";
import "./App.css";

const API_URL = "http://localhost:8080";
const TOKEN_STORAGE_KEY = "authToken";

interface LoginResponse {
  token: string;
}

interface MeResponse {
  email: string;
}

function App() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [token, setToken] = useState(() => localStorage.getItem(TOKEN_STORAGE_KEY) ?? "");
  const [userEmail, setUserEmail] = useState("");
  const [error, setError] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [isLoadingProfile, setIsLoadingProfile] = useState(false);

  const isAuthenticated = token !== "";

  useEffect(() => {
    if (token === "") {
      return;
    }

    void fetchMe(token);
  }, [token]);

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();

    setError("");
    setUserEmail("");
    setIsSubmitting(true);

    try {
      const response = await fetch(`${API_URL}/login`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password }),
      });

      if (!response.ok) {
        setError("Email o password non valide");
        return;
      }

      const data: LoginResponse = await response.json();
      localStorage.setItem(TOKEN_STORAGE_KEY, data.token);
      setToken(data.token);
      setPassword("");
    } catch {
      setError("Backend non raggiungibile");
    } finally {
      setIsSubmitting(false);
    }
  }

  async function fetchMe(authToken: string) {
    setError("");
    setIsLoadingProfile(true);

    try {
      const response = await fetch(`${API_URL}/me`, {
        headers: {
          Authorization: `Bearer ${authToken}`,
        },
      });

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

            <button type="submit" className="btn btn-primary" disabled={isSubmitting}>
              {isSubmitting ? "Accesso..." : "Login"}
            </button>
          </form>
        ) : (
          <div className="protected-view">
            <div className="alert alert-success mb-0">Login riuscito</div>

            <div className="profile-box">
              <span>Utente autenticato</span>
              <strong>{isLoadingProfile ? "Caricamento..." : userEmail}</strong>
            </div>

            <button type="button" className="btn btn-outline-secondary" onClick={handleLogout}>
              Logout
            </button>
          </div>
        )}

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
