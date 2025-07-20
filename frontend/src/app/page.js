import Navbar from "./components/navbar";
import Sidebar from "./components/sidebar";

export default function Home() {
  return (
    <div className="min-h-screen flex flex-col">
      <Navbar />
      <div className="flex flex-1">
        <Sidebar />
        <main className="flex-1 p-4">
        </main>
      </div>
      <footer className="p-4 bg-gray-900">Footer content</footer>
    </div>
  );
}
