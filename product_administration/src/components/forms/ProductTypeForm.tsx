import React from 'react';
import type { Category, SubCategory, ProductType, StoringCondition } from '../../types/Types';
import FeedbackMessage from '../ui/FeedbackMessage';

interface ProductTypeFormProps {
    productTypes: ProductType[];
    categories: Category[];
    subCategories: SubCategory[];
    storingConditions: StoringCondition[];
    selectedId: number | null;
    
    name: string;
    description: string;
    categoryId: number | '';
    subCategoryId: number | '';
    storingConditionId: number | '';
    
    loading: boolean;
    error: string | null;
    successMsg: string | null;
    
    setName: (val: string) => void;
    setDescription: (val: string) => void;
    setCategoryId: (val: number | '') => void;
    setSubCategoryId: (val: number | '') => void;
    setStoringConditionId: (val: number | '') => void;
    setSelectedId: (id: number | null) => void;
    
    handleSubmit: (e: React.FormEvent) => void;
    handleDelete: () => void;
}

export default function ProductTypeForm({ 
    productTypes, categories, subCategories, storingConditions, selectedId, 
    name, description, categoryId, subCategoryId, storingConditionId,
    loading, error, successMsg, 
    setName, setDescription, setCategoryId, setSubCategoryId, setStoringConditionId, setSelectedId, 
    handleSubmit, handleDelete 
}: ProductTypeFormProps) {

    const handleReset = () => {
        setSelectedId(null);
        setName("");
        setDescription("");
        setCategoryId("");
        setSubCategoryId("");
        setStoringConditionId("");
    };

    const handleTypeChange = (val: string) => {
        const id = Number(val);
        if (id === 0) {
            handleReset();
        } else {
            const type = productTypes.find(t => Number(t.id) === id);
            if (type) {
                setSelectedId(id);
                setName(type.name);
                setDescription(type.description);
                setStoringConditionId(Number(type.storing_condition_id));
                
                const sub = subCategories.find(s => Number(s.id) === Number(type.sub_id));
                if (sub) {
                    setCategoryId(Number(sub.category_id));
                    setSubCategoryId(Number(sub.id));
                } else {
                    setCategoryId("");
                    setSubCategoryId("");
                }
            }
        }
    };

    const filteredSubCategories = categoryId 
        ? subCategories.filter(sc => Number(sc.category_id) === Number(categoryId)) 
        : [];

    return (
        <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-4 md:p-8 w-full max-w-3xl">
          <div className="mb-8 flex justify-between items-start">
            <div>
                <h1 className="text-2xl font-bold text-slate-900">
                {selectedId ? 'Terméktípus szerkesztése' : 'Új terméktípus'}
                </h1>
                <p className="text-gray-500 mt-1 text-sm">
                {selectedId 
                    ? 'Szerkessze a kiválasztott terméktípust vagy törölje.' 
                    : 'Adja meg az adatokat az új terméktípus létrehozásához.'}
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
                  onChange={(e) => handleTypeChange(e.target.value)}
                  value={selectedId || 0}
                >
                  <option value={0}>-- Új létrehozása --</option>
                  {productTypes.map(type => (
                    <option key={type.id} value={type.id}>{type.name}</option>
                  ))}
                </select>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="flex flex-col gap-2">
                    <label className="text-xs font-bold text-slate-700 uppercase tracking-wide">
                      Főkategória
                    </label>
                    <select
                      className={`appearance-none w-full border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-3 outline-none transition-all ${
                          selectedId ? 'bg-gray-200 cursor-not-allowed opacity-75' : 'bg-gray-50'
                      }`}
                      value={categoryId}
                      onChange={(e) => {
                          setCategoryId(Number(e.target.value));
                          setSubCategoryId("");
                      }}
                      disabled={!!selectedId}
                    >
                        <option value="" disabled>Válasszon...</option>
                        {categories.map(cat => (
                            <option key={cat.id} value={cat.id}>{cat.name}</option>
                        ))}
                    </select>
                  </div>

                  <div className="flex flex-col gap-2">
                    <label className="text-xs font-bold text-slate-700 uppercase tracking-wide">
                      Alkategória
                    </label>
                    <select
                      className={`appearance-none w-full border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-3 outline-none transition-all ${
                          (!categoryId || !!selectedId) ? 'bg-gray-200 cursor-not-allowed opacity-75' : 'bg-gray-50'
                      }`}
                      value={subCategoryId}
                      onChange={(e) => setSubCategoryId(Number(e.target.value))}
                      disabled={!categoryId || !!selectedId}
                    >
                        <option value="" disabled>Válasszon...</option>
                        {filteredSubCategories.map(sub => (
                            <option key={sub.id} value={sub.id}>{sub.name}</option>
                        ))}
                    </select>
                  </div>
              </div>

              <div className="flex flex-col gap-2">
                <label className="text-xs font-bold text-slate-700 uppercase tracking-wide">
                  Terméktípus neve
                </label>
                <input 
                  type="text" 
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="pl. Teljes kiőrlésű kenyér"
                  className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-2 focus:ring-blue-100 focus:border-blue-500 block p-3 outline-none transition-all"
                />
              </div>

              <div className="flex flex-col gap-2">
                <label className="text-xs font-bold text-slate-700 uppercase tracking-wide">
                  Leírás
                </label>
                <textarea 
                  value={description}
                  onChange={(e) => setDescription(e.target.value)}
                  placeholder="Rövid leírás a típusról..."
                  rows={3}
                  className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-2 focus:ring-blue-100 focus:border-blue-500 block p-3 outline-none transition-all resize-none"
                />
              </div>

              <div className="flex flex-col gap-2">
                <label className="text-xs font-bold text-slate-700 uppercase tracking-wide">
                  Tárolási körülmény
                </label>
                <select 
                  className="w-full bg-gray-50 border border-gray-200 text-gray-700 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-2.5 outline-none transition-all appearance-none"
                  onChange={(e) => setStoringConditionId(Number(e.target.value))}
                  value={storingConditionId || ''}
                >
                  <option value="" disabled>Válasszon...</option>
                  {storingConditions.map(c => (
                    <option key={c.id} value={c.id}>{c.description}</option>
                  ))}
                </select>
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
                disabled={loading || !name.trim() || !subCategoryId || !storingConditionId}
                className="w-full md:w-auto px-6 py-2.5 text-sm font-bold text-white bg-blue-600 rounded-lg hover:bg-blue-700 transition-all shadow-sm shadow-blue-200 focus:outline-none focus:ring-4 focus:ring-blue-100 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {loading ? 'Mentés...' : (selectedId ? 'Módosítás' : 'Létrehozás')}
              </button>
            </div>
          </form>
        </div>
    );
}