import { LayoutDashboard, ShoppingBag, Layers, Package, Tag, Award, Search, Archive } from 'lucide-react';
import { NavLink } from 'react-router-dom';

const SidebarItem = ({ icon: Icon, label, to }: { icon: any, label: string, to: string }) => (
  <NavLink
    to={to}
    className={({ isActive }) => `
      flex items-center gap-3 px-4 py-3 mb-1 cursor-pointer rounded-r-full transition-colors
      ${isActive 
        ? 'bg-blue-100 text-blue-700 border-l-4 border-blue-600' 
        : 'text-gray-500 hover:bg-gray-50 hover:text-gray-900 border-l-4 border-transparent'
      }
    `}
  >
    <Icon size={20} />
    <span className="font-medium text-sm">{label}</span>
  </NavLink>
);

export default function Sidebar() {
    return (
        <aside className="w-64 bg-white border-r border-gray-100 flex flex-col justify-between fixed h-full z-10">
            <div className="pt-6 pr-4">
                <div className="px-6 mb-8 flex items-center gap-2">
                    <div className="w-8 h-8 rounded-lg flex items-center justify-center text-blue-600 font-bold">
                        <svg className="h-6 w-6" fill="none" height="24" stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" viewBox="0 0 24 24" width="24" xmlns="http://www.w3.org/2000/svg"><path d="M12 2L2 7v10l10 5 10-5V7L12 2z"></path><path d="M2 17l10 5"></path><path d="M22 17l-10 5"></path><path d="M12 12l10-5"></path><path d="M12 12v10"></path><path d="M12 12L2 7"></path></svg>
                    </div>
                    <span className="font-bold text-lg text-slate-800">NetStore</span>
                </div>

                <nav>
                    <SidebarItem icon={LayoutDashboard} label="Irányítópult" to="/" />
                    <SidebarItem icon={ShoppingBag} label="Kategóriák" to="/categories" />
                    <SidebarItem icon={Layers} label="Alkategóriák" to="/subcategories" />
                    <SidebarItem icon={Tag} label="Terméktípusok" to="/product-types" />
                    <SidebarItem icon={Package} label="Termékek" to="/products" />
                    <SidebarItem icon={Award} label="Gyártók" to="/brands" />
                    <SidebarItem icon={Archive} label="Tárolási körülmények" to="/storing-condition" />
                    <SidebarItem icon={Search} label="Részletes keresés" to="/search" />
                </nav>
            </div>

            <div className="p-6 border-t border-gray-100">
                <div className="flex items-center gap-3">
                    <div className="w-10 h-10 rounded-full flex items-center justify-center bg-gray-100">
                        <svg
                            xmlns="http://www.w3.org/2000/svg"
                            className="aspect-square w-8 h-8 rounded-full bg-cover bg-center profile-icon fill-blue-600"
                            viewBox="0 0 640 640"
                        >
                            <path
                                d="M320 312C386.3 312 440 258.3 440 192C440 125.7 386.3 72 320 72C253.7 72 200 125.7 200 192C200 258.3 253.7 312 320 312zM290.3 368C191.8 368 112 447.8 112 546.3C112 562.7 125.3 576 141.7 576L498.3 576C514.7 576 528 562.7 528 546.3C528 447.8 448.2 368 349.7 368L290.3 368z"
                            />
                    </svg>
                    </div>
                    <div>
                        <p className="text-sm font-bold text-slate-900">Admin</p>
                        <p className="text-xs text-gray-500 cursor-pointer hover:text-blue-600">Kijelentkezés</p>
                    </div>
                </div>
            </div>
        </aside>
    );
}