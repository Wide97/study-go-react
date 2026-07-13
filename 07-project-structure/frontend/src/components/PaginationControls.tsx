import type { PaginationControlsProps } from "../types";

function PaginationControls({
  page,
  totalPages,
  onPrevious,
  onNext,
}: PaginationControlsProps) {
  return (
    <div className="d-flex gap-2 mt-3">
      <button
        type="button"
        className="btn btn-outline-secondary"
        onClick={onPrevious}
        // Prima pagina: non ha senso andare indietro.
        disabled={page === 1}
      >
        Precedente
      </button>

      <span className="align-self-center">
        Pagina {page} di {totalPages}
      </span>

      <button
        type="button"
        className="btn btn-outline-secondary"
        onClick={onNext}
        // Ultima pagina: non ha senso andare avanti.
        disabled={page >= totalPages}
      >
        Successiva
      </button>
    </div>
  );
}

export { PaginationControls };
