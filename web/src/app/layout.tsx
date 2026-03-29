import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Qube",
  description: "A social network where you never miss a post.",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ja">
      <body>{children}</body>
    </html>
  );
}
