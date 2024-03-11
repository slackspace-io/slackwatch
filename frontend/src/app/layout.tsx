import type { Metadata } from "next";
import { Inter } from "next/font/google";
import "./globals.css";

const inter = Inter({ subsets: ["latin"] });
export const metadata: Metadata = {
  title: "SlackWatch - Kubernetes Container Monitoring",
  description: "Monitor your Kubernetes containers and stay ahead with SlackWatch.",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <head>
        {/* Metadata can be expanded here for SEO purposes */}
        <title>SlackWatch - Kubernetes Container Monitoring</title>
        <meta name="description" content="SlackWatch is your go-to tool for monitoring containers within your Kubernetes cluster, ensuring your deployments are up-to-date and secure." />
      </head>
      <body className={`${inter.className} bg-gray-50`}>{children}</body>
    </html>
  );
}
