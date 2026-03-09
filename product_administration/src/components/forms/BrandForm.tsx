import type React from "react";
import type { Brand } from "../../types/Types"
import FeedbackMessage from "../ui/FeedbackMessage";

interface BrandFormProps {
    brands: Brand[];
    selectedId: number | null;
    name: string;
    isOwn: boolean;
    isTemporary: boolean;
    loading: boolean;
    error: string | null;
    successMsg: string | null;
    setName: (val: string) => void;
    setIsOwn: (val: boolean) => void;
    setIsTemporary: (val: boolean) => void;
    setSelectedId: (id: number | null) => void;
    handleSubmit: (e: React.FormEvent) => void;
    handleDelete: () => void;
}

export default function BrandForm({brands, selectedId, name, isOwn, isTemporary, loading, error, successMsg, setName, setIsOwn, setIsTemporary, setSelectedId, handleSubmit, handleDelete}: BrandFormProps) {
    const handleReset = () => {
        setSelectedId(null);
        setName("");
        setIsOwn(false);
        setIsTemporary(false);
    };

    return(
        <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-4 md:p-8 w-full max-w-3xl">
          <div className="mb-8 flex justify-between items-start">
            <div>
                <h1 className="text-2xl font-bold text-slate-900">
                {selectedId ? 'Márka szerkesztése' : 'Új márka'}
                </h1>
                <p className="text-gray-500 mt-1 text-sm">
                {selectedId 
                    ? 'Szerkessze a kiválasztott márkát vagy törölje.' 
                    : 'Adjon hozzá új márkát.'}
                </p>
            </div>
            {selectedId && (
                <button 
                    onClick={handleReset}
                    className="text-xs font-medium text-blue-600 hover:text-blue-800 bg-blue-50 px-3 py-1.5 rounded-md transition-colors"
                >
                    + Új felvétele
                </button>
            )}
          </div>

          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="grid grid-cols-1 gap-6">
              <div className="flex flex-col gap-2">
                <label className="text-xs font-bold text-slate-700 uppercase tracking-wide">
                  Szerkesztendő kiválasztása
                </label>
                <select 
                  className="w-full bg-gray-50 border border-gray-200 text-gray-700 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-2.5 outline-none transition-all appearance-none"
                  onChange={(e) => {
                    const id = Number(e.target.value);
                    if (id === 0) {
                        handleReset();
                    } else {
                        const brand = brands.find(b => Number(b.id) === id);
                        if (brand) {
                            setSelectedId(id);
                            setName(brand.name);
                            setIsOwn(Number(brand.is_own) === 1);
                            setIsTemporary(Number(brand.is_temporary) === 1);
                        }
                    }
                  }}
                  value={selectedId || 0}
                >
                  <option value={0}>-- Új létrehozása --</option>
                  {brands.map(b => (
                    <option key={b.id} value={b.id}>{b.name}</option>
                  ))}
                </select>
              </div>

              <div className="flex flex-col gap-2">
                <label className="text-xs font-bold text-slate-700 uppercase tracking-wide">
                  Márka neve
                </label>
                <input 
                  type="text" 
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="pl. Coca-Cola"
                  className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-2 focus:ring-blue-100 focus:border-blue-500 block p-3 outline-none transition-all"
                />
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="flex items-center p-3 border border-gray-200 rounded-lg bg-gray-50">
                    <input 
                        id="is_own"
                        type="checkbox"
                        checked={isOwn}
                        onChange={(e) => setIsOwn(e.target.checked)}
                        className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500"
                    />
                    <label htmlFor="is_own" className="ml-2 text-sm font-medium text-gray-900">
                        Saját márka
                    </label>
                  </div>

                  <div className="flex items-center p-3 border border-gray-200 rounded-lg bg-gray-50">
                    <input 
                        id="is_temporary"
                        type="checkbox"
                        checked={isTemporary}
                        onChange={(e) => setIsTemporary(e.target.checked)}
                        className="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500"
                    />
                    <label htmlFor="is_temporary" className="ml-2 text-sm font-medium text-gray-900">
                        Ideiglenes választék
                    </label>
                  </div>
              </div>
            </div>

            {error && <FeedbackMessage type="error" message={error} />}
            {successMsg && <FeedbackMessage type="success" message={successMsg} />}

            <div className="flex flex-col-reverse md:flex-row items-center justify-end gap-3 md:gap-4 pt-6 mt-2 border-t border-gray-50">
              {selectedId && (
                <button
                  type="button"
                  onClick={handleDelete}
                  disabled={loading}
                  className="w-full md:w-auto px-5 py-2.5 text-sm font-bold text-red-600 bg-red-50 rounded-lg hover:bg-red-100 transition-colors focus:outline-none focus:ring-2 focus:ring-red-100"
                >
                  Törlés
                </button>
              )}

              <button
                type="submit"
                disabled={loading || !name.trim()}
                className="w-full md:w-auto px-6 py-2.5 text-sm font-bold text-white bg-blue-600 rounded-lg hover:bg-blue-700 transition-all shadow-sm shadow-blue-200 focus:outline-none focus:ring-4 focus:ring-blue-100 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {loading ? 'Mentés...' : (selectedId ? 'Módosítás' : 'Létrehozás')}
              </button>
            </div>
          </form>
        </div>
    );
}