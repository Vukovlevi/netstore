import { useAuth, ROLES } from '../context/AuthContext';
import { ShoppingBag, Layers, Tag, Package, Award, Archive, Search } from 'lucide-react';
import { Link } from 'react-router-dom';

export default function Dashboard() {
  const { user, role, canWrite } = useAuth();

  const getRoleDescription = (role: string | null) => {
    switch (role) {
      case ROLES.UZLETVEZETO:
        return 'Teljes hozzáférés az összes funkcióhoz';
      case ROLES.HR:
        return 'Felhasználók és szerződések kezelése';
      case ROLES.RAKTARVEZETO:
        return 'Kategóriák, típusok, gyártók és termékek kezelése';
      case ROLES.RAKTARKEZELO:
        return 'Termékek kezelése';
      case ROLES.PENZTAROS:
        return 'Termékek keresése és mennyiség csökkentése';
      case ROLES.EGYEB:
        return 'Termékek keresése';
      default:
        return 'Ismeretlen jogosultság';
    }
  };

  const quickLinks = [
    { icon: ShoppingBag, label: 'Kategóriák', to: '/categories', resource: 'category' },
    { icon: Layers, label: 'Alkategóriák', to: '/subcategories', resource: 'sub_category' },
    { icon: Tag, label: 'Terméktípusok', to: '/product-types', resource: 'product_type' },
    { icon: Package, label: 'Termékek', to: '/products', resource: 'product' },
    { icon: Award, label: 'Gyártók', to: '/brands', resource: 'brand' },
    { icon: Archive, label: 'Tárolási körülmények', to: '/storing-condition', resource: 'storing_condition' },
    { icon: Search, label: 'Részletes keresés', to: '/search', resource: 'search' },
  ];

  const availableLinks = quickLinks.filter(
    (link) => link.resource === 'search' || canWrite(link.resource)
  );

  return (
    <div className="space-y-8">
      <div>
        <h1 className="text-2xl font-bold text-slate-900 mb-2">Irányítópult</h1>
        <p className="text-gray-500">Üdvözöljük az adminisztrációs felületen.</p>
      </div>

      <div className="bg-gradient-to-r from-blue-500 to-blue-600 rounded-2xl p-4 md:p-6 text-white shadow-lg">
        <div className="flex items-center gap-4">
          <div className="w-16 h-16 bg-white/20 rounded-full flex items-center justify-center">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="w-10 h-10 fill-white"
              viewBox="0 0 640 640"
            >
              <path d="M320 312C386.3 312 440 258.3 440 192C440 125.7 386.3 72 320 72C253.7 72 200 125.7 200 192C200 258.3 253.7 312 320 312zM290.3 368C191.8 368 112 447.8 112 546.3C112 562.7 125.3 576 141.7 576L498.3 576C514.7 576 528 562.7 528 546.3C528 447.8 448.2 368 349.7 368L290.3 368z" />
            </svg>
          </div>
          <div>
            <h2 className="text-xl font-bold">{user?.username || 'Felhasználó'}</h2>
            <p className="text-blue-100 text-sm">{role || 'Ismeretlen szerepkör'}</p>
            <p className="text-blue-200 text-xs mt-1">{getRoleDescription(role)}</p>
          </div>
        </div>
      </div>

      <div>
        <h2 className="text-lg font-semibold text-slate-800 mb-4">Gyors elérés</h2>
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
          {availableLinks.map((link) => (
            <Link
              key={link.to}
              to={link.to}
              className="flex flex-col items-center gap-3 p-4 md:p-6 bg-white rounded-xl border border-gray-100 hover:border-blue-200 hover:shadow-md transition-all group"
            >
              <div className="p-3 bg-gray-50 rounded-lg group-hover:bg-blue-50 transition-colors">
                <link.icon className="w-6 h-6 text-gray-500 group-hover:text-blue-600 transition-colors" />
              </div>
              <span className="text-sm font-medium text-gray-700 group-hover:text-blue-600 transition-colors">
                {link.label}
              </span>
            </Link>
          ))}
        </div>
      </div>

      <div className="bg-gray-50 rounded-xl p-6">
        <h2 className="text-lg font-semibold text-slate-800 mb-4">Jogosultságok</h2>
        <div className="space-y-3">
          <div className="flex items-center gap-3">
            <div className={`w-3 h-3 rounded-full ${canWrite('category') ? 'bg-green-500' : 'bg-gray-300'}`} />
            <span className="text-sm text-gray-600">Kategóriák, alkategóriák, típusok, gyártók, tárolási körülmények kezelése</span>
          </div>
          <div className="flex items-center gap-3">
            <div className={`w-3 h-3 rounded-full ${canWrite('product') ? 'bg-green-500' : 'bg-gray-300'}`} />
            <span className="text-sm text-gray-600">Termékek kezelése</span>
          </div>
          <div className="flex items-center gap-3">
            <div className="w-3 h-3 rounded-full bg-green-500" />
            <span className="text-sm text-gray-600">Termékek keresése</span>
          </div>
        </div>
      </div>

      <a
        href="http://localhost:8000"
        className="inline-flex px-6 py-3 text-sm font-bold text-white bg-blue-600 rounded-lg hover:bg-blue-700 transition-all shadow-sm shadow-blue-200 focus:outline-none focus:ring-4 focus:ring-blue-100"
      >
        Ugrás a központi felületre
      </a>
    </div>
  );
}
