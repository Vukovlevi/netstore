import { LayoutDashboard, ShoppingBag, ClipboardList, Users, UserCircle } from 'lucide-react';

const SidebarItem = ({ icon: Icon, label, active = false }: { icon: any, label: string, active?: boolean }) => (
  <div className={`flex items-center gap-3 px-4 py-3 mb-1 cursor-pointer rounded-r-full transition-colors ${
    active 
      ? 'bg-blue-100 text-blue-700 border-l-4 border-blue-600' 
      : 'text-gray-500 hover:bg-gray-50 hover:text-gray-900'
  }`}>
    <Icon size={20} />
    <span className="font-medium text-sm">{label}</span>
  </div>
);

export default function Sidebar() {
    return (
        <aside className="w-64 bg-white border-r border-gray-100 flex flex-col justify-between fixed h-full z-10">
            <div className="pt-6 pr-4">
                <div className="px-6 mb-8 flex items-center gap-2">
                    <div className="w-8 h-8 bg-blue-600 rounded-lg flex items-center justify-center text-white font-bold">
                        <ShoppingBag size={18} />
                    </div>
                    <span className="font-bold text-lg text-slate-800">Netstore</span>
                </div>

                <nav>
                    <SidebarItem icon={LayoutDashboard} label="Irányítópult" />
                    <SidebarItem icon={ShoppingBag} label="Termékek / Kategóriák" active />
                    <SidebarItem icon={ClipboardList} label="Rendelések" />
                    <SidebarItem icon={Users} label="Dolgozók" />
                </nav>
            </div>

            <div className="p-6 border-t border-gray-100">
                <div className="flex items-center gap-3">
                    <div className="w-10 h-10 rounded-full bg-orange-100 flex items-center justify-center">
                        <UserCircle className="text-orange-500" />
                    </div>
                    <div>
                        <p className="text-sm font-bold text-slate-900">Admin Felhasználó</p>
                        <p className="text-xs text-gray-500 cursor-pointer hover:text-blue-600">Kijelentkezés</p>
                    </div>
                </div>
            </div>
        </aside>
    );
}