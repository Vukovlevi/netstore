import { ShieldX } from 'lucide-react';
import { useAuth } from '../context/AuthContext';
import { Link } from 'react-router-dom';

interface AccessDeniedProps {
  resource: string;
  requiredRoles?: string[];
}

export default function AccessDenied({ resource, requiredRoles }: AccessDeniedProps) {
  const { role } = useAuth();

  const resourceNames: Record<string, string> = {
    category: 'Kategóriák',
    sub_category: 'Alkategóriák',
    product_type: 'Terméktípusok',
    product: 'Termékek',
    brand: 'Gyártók',
    storing_condition: 'Tárolási körülmények',
    ingredient: 'Összetevők',
    user: 'Felhasználók',
    contract: 'Szerződések',
    contract_type: 'Szerződéstípusok',
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-[60vh] text-center px-4">
      <div className="bg-red-50 rounded-full p-6 mb-6">
        <ShieldX className="w-16 h-16 text-red-500" />
      </div>
      <h1 className="text-2xl font-bold text-gray-900 mb-2">
        Hozzáférés megtagadva
      </h1>
      <p className="text-gray-600 mb-4 max-w-md">
        Önnek nincs jogosultsága a(z) <span className="font-semibold">{resourceNames[resource] || resource}</span> oldal megtekintéséhez.
      </p>
      {role && (
        <p className="text-sm text-gray-500 mb-4">
          Az Ön jelenlegi szerepköre: <span className="font-medium text-gray-700">{role}</span>
        </p>
      )}
      {requiredRoles && requiredRoles.length > 0 && (
        <div className="bg-gray-50 rounded-lg p-4 max-w-md mb-6">
          <p className="text-sm text-gray-600 mb-2">Ehhez az oldalhoz szükséges szerepkör:</p>
          <div className="flex flex-wrap gap-2 justify-center">
            {requiredRoles.map((r) => (
              <span
                key={r}
                className="px-3 py-1 bg-blue-100 text-blue-700 text-sm rounded-full"
              >
                {r}
              </span>
            ))}
          </div>
        </div>
      )}
      <Link
        to="/"
        className="px-6 py-3 text-sm font-bold text-white bg-blue-600 rounded-lg hover:bg-blue-700 transition-all"
      >
        Vissza az irányítópultra
      </Link>
    </div>
  );
}
