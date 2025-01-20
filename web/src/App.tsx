import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router';
import './styles/global.css';
import { Login } from './pages/Login/Login';
import { Browse } from './pages/Browse/Browse';
import { NotFound } from './pages/404';

export function App() {
  return (
    <Router>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/browse/*" element={<Browse />} />
        <Route path="/404" element={<NotFound />} />
        <Route path="*" element={<Navigate to="/browse" replace />} />
      </Routes>
    </Router>
  );
}
