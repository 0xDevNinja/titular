"use client";

export default function AgentError({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <main
      style={{ maxWidth: 800, margin: "0 auto", padding: "40px 24px", textAlign: "center" }}
      role='alert'
    >
      <h1 style={{ color: "#dc2626" }}>Something went wrong</h1>
      <p style={{ color: "#6b7280", marginBottom: 24 }}>
        {error.message || "Failed to load agent page."}
      </p>
      <button
        type='button'
        onClick={reset}
        style={{
          padding: "10px 24px",
          background: "#2563eb",
          color: "#fff",
          border: "none",
          borderRadius: 6,
          cursor: "pointer",
        }}
      >
        Try again
      </button>
    </main>
  );
}
