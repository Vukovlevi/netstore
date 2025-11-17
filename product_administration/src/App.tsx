import Sidebar from './components/layout/Sidebar';
import CategoryManagement from './views/CategoryManagement';

export default function App() {
  return (
    <div className="flex min-h-screen bg-[#f8f9fb] font-sans text-slate-800">
      <Sidebar />
      <CategoryManagement />
    </div>
  );
}