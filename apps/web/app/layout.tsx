import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "Titular",
  description: "Multi-chain agent launchpad",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang='en'>
      <body>{children}</body>
    </html>
  );
}
