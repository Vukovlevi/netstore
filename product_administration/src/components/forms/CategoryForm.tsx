import React from 'react';
import type { Category } from '../../types/Types';
import FeedbackMessage from '../ui/FeedbackMessage';

interface CategoryFormProps {
    categories: Category[];
    selectedId: number | null;
    name: string;
    loading: boolean;
    error: string | null;
    successMsg: string | null;
    setName: (name: string) => void;
    setSelectedId: (id: number | null) => void;
    handleSubmit: (e: React.FormEvent) => void;
    handleDelete: () => void;
}

export default function CategoryForm({ 
    categories, selectedId, name, loading, error, successMsg, 
    setName, setSelectedId, handleSubmit, handleDelete 
}: CategoryFormProps) {

    const handleReset = () => {
        setSelectedId(null);
        setName("");
    };

    return (
        <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-4 md:p-8 w-full max-w-3xl">
          
          <div className="mb-8 flex justify-between items-start">
            <div>
                <h1 className="text-2xl font-bold text-slate-900">
                {selectedId ? 'Kategória szerkesztése' : 'Új kategória'}
                </h1>
                <p className="text-gray-500 mt-1 text-sm">
                {selectedId 
                    ? 'Szerkessze a kiválasztott kategóriát vagy törölje.' 
                    : 'Adja meg az új kategória nevét a létrehozáshoz.'}
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
                        const cat = categories.find(c => c.id === id);
                        setSelectedId(id);
                        setName(cat ? cat.name : "");
                    }
                  }}
                  value={selectedId || 0}
                >
                  <option value={0}>-- Új létrehozása --</option>
                  {categories.map(cat => (
                    <option key={cat.id} value={cat.id}>{cat.name}</option>
                  ))}
                </select>
              </div>

              <div className="flex flex-col gap-2">
                <label className="text-xs font-bold text-slate-700 uppercase tracking-wide">
                  Kategória neve
                </label>
                <input 
                  type="text" 
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="pl. Prémium alma"
                  className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-2 focus:ring-blue-100 focus:border-blue-500 block p-3 outline-none transition-all"
                />
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