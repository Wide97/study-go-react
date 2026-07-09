import { useState } from "react";
import "./App.css";

interface LoginResponse {
  token: string;
}

interface MeResponse {
  email: string;
}

function App() {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [token, setToken] = useState("");
  const [userEmail, setUserEmail] = useState("");

  async function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault();

    const response = await fetch("http://localhost:8080/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email, password }),
    });

    if (!response.ok) {
      console.error("Login fallito");
      return;
    }

    const data: LoginResponse = await response.json();
    setToken(data.token);
    console.log(data.token);
  }

  async function fetchMe() {
    const response = await fetch("http://localhost:8080/me", {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    });

    if (!response.ok) {
      console.error("Chiamata /me fallita");
      return;
    }

    const data: MeResponse = await response.json();
    setUserEmail(data.email);
  }

  return (
    <main className="app-shell">
      <section className="app-panel">
        <h1>Login</h1>

        <form onSubmit={handleSubmit}>
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
            />
          </div>

          <button type="submit" className="btn btn-primary">
            Login
          </button>
        </form>

        {token !== "" && (
          <>
            <div className="alert alert-success mt-3">Login riuscito</div>

            <button
              type="button"
              className="btn btn-outline-primary mt-2"
              onClick={fetchMe}
            >
              Carica profilo
            </button>
          </>
        )}
        {userEmail !== "" && (
          <p className="mt-3 mb-0">
            Utente autenticato: <strong>{userEmail}</strong>
          </p>
        )}
      </section>
    </main>
  );
}

export default App;
