import { useEffect, useState } from 'react'
import './App.css'

function App() {

  // useState restituisce due cose: il valore attuale dello stato (todos)
  // e una funzione per aggiornarlo (setTodos). Quando setTodos viene chiamata,
  // React ridisegna il componente con il nuovo valore.
  // <Todo[]> è il tipo generico: dice a TypeScript che questo stato
  // conterrà sempre un array di Todo (anche se parte vuoto []).
  const [todos, setTodos] = useState<Todo[]>([]);

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

  return (
    <>
    <div>
      <h1>Todo List</h1>
      <ul>
        {/* todos.map trasforma ogni elemento dell'array in un elemento JSX.
            key={todo.id} è obbligatorio con le liste in React: aiuta a capire
            quale elemento è cambiato/aggiunto/rimosso senza ridisegnare tutto. */}
        {todos.map(todo => (
          <li key={todo.id}>
            <span>{todo.title}</span>
            <span>{todo.done ? ' (Done)' : ' (Not Done)'}</span>
          </li>
        ))}
      </ul>
    </div>

    </>
  )
}

export default App
