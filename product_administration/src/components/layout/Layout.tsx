import { Outlet } from 'react-router-dom';
import Sidebar from './Sidebar';

export default function Layout() {
    return (
        <div className="flex min-h-screen bg-[#f8f9fb] font-sans text-slate-800">
            <Sidebar />
            <div className="flex-1 ml-64 p-12">
                <Outlet />
            </div>
        </div>
    );
}