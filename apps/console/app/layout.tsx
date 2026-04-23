import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Titular Console",
  description: "Agent management console",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
