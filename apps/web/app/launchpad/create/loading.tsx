export default function Loading() {
  return (
    <main style={{ maxWidth: 680, margin: "0 auto", padding: "40px 24px" }}>
      <p aria-live='polite' aria-busy='true'>
        Loading wizard...
      </p>
    </main>
  );
}
