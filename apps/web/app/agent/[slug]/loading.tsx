export default function AgentLoading() {
  return (
    <main
      style={{ maxWidth: 800, margin: "0 auto", padding: "40px 24px" }}
      aria-busy='true'
      aria-label='Loading agent page'
    >
      <div
        style={{
          height: 32,
          width: 240,
          background: "#e2e8f0",
          borderRadius: 6,
          marginBottom: 24,
          animation: "pulse 1.5s ease-in-out infinite",
        }}
      />
      <div
        style={{
          height: 200,
          background: "#e2e8f0",
          borderRadius: 8,
          marginBottom: 24,
        }}
      />
      <div
        style={{
          height: 300,
          background: "#e2e8f0",
          borderRadius: 8,
        }}
      />
    </main>
  );
}
