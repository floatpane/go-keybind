import type { Metadata } from "next";
import "./globals.css";
import "highlight.js/styles/github-dark.css";

export const metadata: Metadata = {
	title: "go-keybind",
	description:
		"Key-binding config for Go TUI apps. Parse and normalize key strings, load/save JSON keybind files with merge-with-defaults semantics, and validate for intra-area conflicts.",
};

export default function RootLayout({
	children,
}: {
	children: React.ReactNode;
}) {
	return (
		<html lang="en">
			<body>
				<header className="site-header">
					<a href="/" className="brand">
						go-keybind
					</a>
					<nav>
						<a href="https://github.com/floatpane/go-keybind">GitHub</a>
					</nav>
				</header>
				<main>{children}</main>
			</body>
		</html>
	);
}

export const viewport = { width: "device-width", initialScale: 1 };
