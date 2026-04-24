"use client";

// biome-ignore lint/suspicious/noShadowRestrictedNames: Next.js App Router requires the default export to be named Error
export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <main style={{ maxWidth: 680, margin: "0 auto", padding: "40px 24px" }}>
      <h2>Something went wrong</h2>
      <p role='alert'>{error.message}</p>
      <button type='button' onClick={reset} style={{ marginTop: 16, padding: "8px 20px" }}>
        Try again
      </button>
    </main>
  );
}
