import { useState } from 'react';
import { Outlet } from 'react-router-dom';
import { Menu } from 'lucide-react';
import Sidebar from './Sidebar';

export default function Layout() {
    const [sidebarOpen, setSidebarOpen] = useState(false);

    return (
        <div className="flex min-h-screen bg-[#f8f9fb] font-sans text-slate-800">
            <Sidebar isOpen={sidebarOpen} onClose={() => setSidebarOpen(false)} />
            <div className="flex-1 md:ml-64 p-4 md:p-8 lg:p-12">
                <button
                    onClick={() => setSidebarOpen(true)}
                    className="md:hidden mb-4 p-2 rounded-lg bg-white border border-gray-200 shadow-sm"
                >
                    <Menu className="w-6 h-6 text-gray-700" />
                </button>
                <Outlet />
            </div>
        </div>
    );
}