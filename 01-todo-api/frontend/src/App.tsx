import { useEffect, useState } from 'react'
import './App.css'

function App() {

  // useState restituisce due cose: il valore attuale dello stato (todos)
  // e una funzione per aggiornarlo (setTodos). Quando setTodos viene chiamata,
  // React ridisegna il componente con il nuovo valore.
  // <Todo[]> è il tipo generico: dice a TypeScript che questo stato
  // conterrà sempre un array di Todo (anche se parte vuoto []).
  const [todos, setTodos] = useState<Todo[]>([]);
  const [title, setTitle] = useState('');

  // Interfaccia TypeScript: descrive la forma dei dati che arrivano
  // dal backend Go (stessi campi della struct Todo lato Go).
  interface Todo {
    id: number;
    title: string;
    done: boolean;
  }

  // useEffect esegue del codice "a lato" del rendering normale
  // (qui: una chiamata di rete). L'array vuoto [] come secondo argomento
  // significa "esegui solo una volta, al primo render del componente"
  // — senza, girerebbe ad ogni render, causando un loop di richieste.
  useEffect(() => {
    // fetch restituisce una Promise e ritorna subito: il codice nei .then
    // gira più avanti, quando la risposta è arrivata (asincrono).
    fetch('http://localhost:8080/todos')
      .then(response => response.json()) // converte il body della risposta in JSON
      .then(data => setTodos(data)) // aggiorna lo stato -> React ridisegna la lista
      .catch(error => console.error('Errore durante il recupero dei todo:', error)); // errori di rete/parsing
  }, []);

  // Funzione normale (NON un useEffect): gira solo quando viene chiamata,
  // cioè quando il form viene inviato (vedi onSubmit più sotto). useEffect
  // serve per reagire al montaggio/cambiamento di stato, non alle azioni utente.
  function handleSubmit(e: React.FormEvent<HTMLFormElement>) {
    e.preventDefault(); // evita il reload della pagina (comportamento di default di un <form> HTML)

    fetch('http://localhost:8080/todos', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title, done: false }), // JSON.stringify: oggetto JS -> stringa JSON
    })
      .then(response => response.json()) // il backend risponde con la todo creata (con id assegnato)
      .then(newTodo => setTodos([...todos, newTodo])) // ...todos copia l'array esistente, poi aggiunge la nuova
      .catch(error => console.error('Errore durante la creazione:', error));

    setTitle(''); // svuota l'input, pronto per la prossima todo
  }

  // Funzione richiamata al click (onClick), non un useEffect: un effect gira
  // al montaggio/al cambiare di certi valori, non "quando l'utente clicca qualcosa".
  // Il backend (toggleTodo in Go) inverte lui stesso Done e ignora il body,
  // quindi qui non serve costruire nulla da mandare: solo il metodo PUT.
  function handleToggleDone(id: number) {

fetch(`http://localhost:8080/todos/${id}`, { method: 'PUT' })
  .then(response => response.json()) // il server risponde con la todo aggiornata (fonte di verità)
  // todos.map scorre tutto l'array e sostituisce SOLO l'elemento con id corrispondente,
  // lasciando gli altri invariati: è così che si aggiorna un singolo elemento
  // di uno stato-array senza mutare l'array originale (React non lo noterebbe).
  .then(updatedTodo => setTodos(todos.map(todo => (todo.id === id ? updatedTodo : todo))))
  .catch(error => console.error('Errore durante l\'aggiornamento:', error));
  }

  // Stesso schema di handleToggleDone, ma il backend risponde 204 No Content
  // (nessun body): niente .then(response => response.json()), andrebbe in errore
  // provando a parsare un body vuoto. todos.filter è il simmetrico "in negativo"
  // di map: tiene tutti gli elementi TRANNE quello con l'id cancellato.
  function handleDelete(id: number) {

  fetch(`http://localhost:8080/todos/${id}`, { method: 'DELETE' })
    .then(() => setTodos(todos.filter(todo => todo.id !== id)))
    .catch(error => console.error('Errore durante la cancellazione:', error));
}

  return (
    // Bootstrap è solo CSS: le classi (container, mb-4, ecc.) applicano stili
    // predefiniti, non c'è nessuna logica React coinvolta. "container" centra
    // il contenuto e gli dà un padding laterale responsive.
    <div className="container py-5" style={{ maxWidth: '600px' }}>
      <h1 className="mb-4 text-center">Todo List</h1>

      {/* Il <form> avvolge input+bottone: onSubmit scatta quando premi
          "Aggiungi" O premi Invio dentro l'input (comportamento nativo HTML).
          "input-group" è la classe Bootstrap che unisce visivamente input e bottone. */}
      <form onSubmit={handleSubmit} className="mb-4">
        <div className="input-group">
          <input
            type="text"
            className="form-control"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Nuova todo..."
          />
          <button type="submit" className="btn btn-primary">Aggiungi</button>
        </div>
      </form>

      <ul className="list-group">
        {/* todos.map trasforma ogni elemento dell'array in un elemento JSX.
            key={todo.id} è obbligatorio con le liste in React: aiuta a capire
            quale elemento è cambiato/aggiunto/rimosso senza ridisegnare tutto. */}
        {todos.map(todo => (
          <li key={todo.id} className="list-group-item d-flex justify-content-between align-items-center">
            {/* className è un'espressione JS: qui una classe diversa in base
                a todo.done, per dare un segnale visivo chiaro (testo barrato)
                quando la todo è completata — non solo testo, come prima. */}
            <span className={todo.done ? 'text-decoration-line-through text-muted' : ''}>
              {todo.title}
            </span>
            <div className="d-flex gap-2 align-items-center">
              {/* Template string (backtick) per comporre la classe del badge
                  dinamicamente in base allo stato, stesso concetto usato
                  negli URL con l'id (es. `.../todos/${id}`). */}
              <span className={`badge ${todo.done ? 'bg-success' : 'bg-secondary'}`}>
                {todo.done ? 'Fatta' : 'Da fare'}
              </span>
              <button className="btn btn-sm btn-outline-secondary" onClick={() => handleToggleDone(todo.id)}>
                Spunta
              </button>
              <button className="btn btn-sm btn-outline-danger" onClick={() => handleDelete(todo.id)}>
                Elimina
              </button>
            </div>
          </li>
        ))}
      </ul>
    </div>
  )
}

export default App
