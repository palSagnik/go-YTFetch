import { BrowserRouter as Router, Route, Routes, Link } from "react-router-dom";
import VideoList from "./components/VideoList";
import { Youtube } from 'lucide-react';

function App() {
  return (
    <Router>
      <div className="flex flex-col min-h-screen bg-[#0f1117]">
        {/* Navbar */}
        <nav className="bg-[#1a1d24] border-b border-gray-800 fixed w-full z-10">
          <div className="max-w-[1800px] mx-auto px-4">
            <div className="flex items-center h-14">
              <Link to="/" className="flex items-center space-x-2">
                <Youtube className="h-6 w-6 text-red-500" />
                <span className="font-semibold text-lg text-white">VideoHub</span>
              </Link>
            </div>
          </div>
        </nav>

        {/* Main Content */}
        <main className="flex-grow pt-14">
          <div className="max-w-[1800px] mx-auto p-4">
            <Routes>
              <Route path="/:cursor?" element={<VideoList />} />
            </Routes>
          </div>
        </main>
      </div>
    </Router>
  );
}

export default App;

