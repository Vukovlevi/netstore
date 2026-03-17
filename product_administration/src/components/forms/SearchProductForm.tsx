import React from "react";
import type {
  Category,
  SubCategory,
  ProductType,
  Brand,
  StoringCondition,
} from "../../types/Types";
import type { SearchFilters } from "../../services/searchFilters";
import { Filter, ChevronRight, ChevronLeft, Layers, Tag, ScanSearch } from "lucide-react";

interface SearchProductFormProps {
  categories: Category[];
  subCategories: SubCategory[];
  productTypes: ProductType[];
  brands: Brand[];
  storingConditions: StoringCondition[];

  filters: SearchFilters;
  setFilters: React.Dispatch<React.SetStateAction<SearchFilters>>;

  activeTab: "relations" | "attributes";
  setActiveTab: (tab: "relations" | "attributes") => void;

  handleSearch: (e?: React.FormEvent) => void;
  handleNetworkSearch: (e?: React.FormEvent) => void;
  clearFilters: () => void;
  loading: "local" | "network" | null;
}

export default function SearchProductForm({
  categories,
  subCategories,
  productTypes,
  brands,
  storingConditions,
  filters,
  setFilters,
  activeTab,
  setActiveTab,
  handleSearch,
  handleNetworkSearch,
  clearFilters,
  loading,
}: SearchProductFormProps) {
  const handleCategoryChange = (value: string) => {
    const catId = value ? Number(value) : undefined;
    setFilters((prev) => ({
      ...prev,
      category_id: catId,
      sub_category_id: undefined,
      type_id: undefined,
      storing_condition_id: undefined,
    }));
  };

  const handleSubCategoryChange = (value: string) => {
    const subId = value ? Number(value) : undefined;
    if (!subId) {
      setFilters((prev) => ({
        ...prev,
        sub_category_id: undefined,
        type_id: undefined,
        storing_condition_id: undefined,
      }));
      return;
    }
    const sub = subCategories.find((s) => Number(s.id) === subId);
    setFilters((prev) => ({
      ...prev,
      sub_category_id: subId,
      category_id: sub ? Number(sub.category_id) : prev.category_id,
      type_id: undefined,
      storing_condition_id: undefined,
    }));
  };

  const handleTypeChange = (value: string) => {
    const typeId = value ? Number(value) : undefined;
    if (!typeId) {
      setFilters((prev) => ({
        ...prev,
        type_id: undefined,
        storing_condition_id: undefined,
      }));
      return;
    }
    const type = productTypes.find((t) => Number(t.id) === typeId);
    if (type) {
      const sub = subCategories.find((s) => Number(s.id) === Number(type.sub_id));
      setFilters((prev) => ({
        ...prev,
        type_id: typeId,
        sub_category_id: sub ? Number(sub.id) : undefined,
        category_id: sub ? Number(sub.category_id) : undefined,
        storing_condition_id: Number(type.storing_condition_id),
      }));
    }
  };

  const handleStoringConditionChange = (value: string) => {
    const condId = value ? Number(value) : undefined;
    setFilters((prev) => ({
      ...prev,
      storing_condition_id: condId,
    }));
  };

  const getAvailableSubCategories = () => {
    if (!subCategories || subCategories.length === 0) return [];
    if (filters.category_id && filters.category_id > 0) {
      return subCategories.filter(
        (sc) => Number(sc.category_id) === Number(filters.category_id)
      );
    }
    return subCategories;
  };

  const getAvailableTypes = () => {
    if (!productTypes || productTypes.length === 0) return [];
    if (filters.sub_category_id && filters.sub_category_id > 0) {
      return productTypes.filter(
        (pt) => Number(pt.sub_id) === Number(filters.sub_category_id)
      );
    }
    if (filters.category_id && filters.category_id > 0) {
      const subIds = subCategories
        .filter((sc) => Number(sc.category_id) === Number(filters.category_id))
        .map((sc) => Number(sc.id));
      return productTypes.filter((pt) => subIds.includes(Number(pt.sub_id)));
    }
    return productTypes;
  };

  const getAvailableStoringConditions = () => {
    if (!storingConditions || storingConditions.length === 0) return [];
    return storingConditions;
  };

  const getActiveFiltersCount = () => {
    let count = 0;
    if (filters.category_id) count++;
    if (filters.sub_category_id) count++;
    if (filters.type_id) count++;
    if (filters.storing_condition_id) count++;
    if (filters.brand_id) count++;
    return count;
  };

  const goToNextTab = () => {
    setActiveTab("attributes");
  };

  const goToPrevTab = () => {
    setActiveTab("relations");
  };

  return (
    <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
      <div className="p-4 md:p-6 border-b border-gray-100 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <h1 className="text-2xl font-bold text-slate-900">Komplex szűrés</h1>
        </div>
        <button
          type="button"
          onClick={clearFilters}
          className="text-sm text-gray-500 hover:text-red-600 flex items-center gap-1 transition-colors"
        >
          Szűrők törlése
        </button>
      </div>

      <div className="flex border-b border-gray-100">
        <button
          type="button"
          className={`flex-1 py-4 text-sm font-medium text-center transition-colors flex items-center justify-center gap-2 ${
            activeTab === "relations"
              ? "text-blue-600 border-b-2 border-blue-600 bg-blue-50/50"
              : "text-gray-500 hover:text-gray-700 hover:bg-gray-50"
          }`}
          onClick={() => setActiveTab("relations")}
        >
          <Layers className="w-4 h-4" />
          <span>1. Kapcsolatok</span>
          {getActiveFiltersCount() > 0 && (
            <span className="bg-blue-600 text-white text-xs px-2 py-0.5 rounded-full">
              {getActiveFiltersCount()}
            </span>
          )}
        </button>
        <button
          type="button"
          className={`flex-1 py-4 text-sm font-medium text-center transition-colors flex items-center justify-center gap-2 ${
            activeTab === "attributes"
              ? "text-blue-600 border-b-2 border-blue-600 bg-blue-50/50"
              : "text-gray-500 hover:text-gray-700 hover:bg-gray-50"
          }`}
          onClick={() => setActiveTab("attributes")}
        >
          <Tag className="w-4 h-4" />
          <span>2. Tulajdonságok</span>
        </button>
      </div>

      <form onSubmit={handleSearch} className="p-4 md:p-6">
        {activeTab === "relations" && (
          <>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label className="text-xs font-bold text-slate-700 uppercase tracking-wide block mb-2">
                  Kategória
                </label>
                <select
                  className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-3"
                  value={filters.category_id || ""}
                  onChange={(e) => handleCategoryChange(e.target.value)}
                >
                  <option value="">Válasszon...</option>
                  {categories.map((c) => (
                    <option key={c.id} value={c.id}>
                      {c.name}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label className="text-xs font-bold text-slate-700 uppercase tracking-wide block mb-2">
                  Alkategória
                </label>
                <select
                  className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-3"
                  value={filters.sub_category_id || ""}
                  onChange={(e) => handleSubCategoryChange(e.target.value)}
                >
                  <option value="">Válasszon...</option>
                  {getAvailableSubCategories().map((s) => (
                    <option key={s.id} value={s.id}>
                      {s.name}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label className="text-xs font-bold text-slate-700 uppercase tracking-wide block mb-2">
                  Terméktípus
                </label>
                <select
                  className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-3"
                  value={filters.type_id || ""}
                  onChange={(e) => handleTypeChange(e.target.value)}
                >
                  <option value="">Válasszon...</option>
                  {getAvailableTypes().map((t) => (
                    <option key={t.id} value={t.id}>
                      {t.name}
                    </option>
                  ))}
                </select>
              </div>

              <div>
                <label className="text-xs font-bold text-slate-700 uppercase tracking-wide block mb-2">
                  Tárolási körülmény
                </label>
                <select
                  className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-3"
                  value={filters.storing_condition_id || ""}
                  onChange={(e) => handleStoringConditionChange(e.target.value)}
                  disabled={!!filters.type_id}
                >
                  <option value="">Válasszon...</option>
                  {getAvailableStoringConditions().map((sc) => (
                    <option key={sc.id} value={sc.id}>
                      {sc.description}
                    </option>
                  ))}
                </select>
              </div>

              <div className="md:col-span-2">
                <label className="text-xs font-bold text-slate-700 uppercase tracking-wide block mb-2">
                  Gyártó
                </label>
                <select
                  className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-3"
                  value={filters.brand_id || ""}
                  onChange={(e) =>
                    setFilters((prev) => ({
                      ...prev,
                      brand_id: e.target.value ? Number(e.target.value) : undefined,
                    }))
                  }
                >
                  <option value="">Válasszon...</option>
                  {brands.map((b) => (
                    <option key={b.id} value={b.id}>
                      {b.name}
                    </option>
                  ))}
                </select>
              </div>
            </div>

            <div className="mt-8 pt-6 border-t border-gray-100 flex justify-end">
              <button
                type="button"
                onClick={goToNextTab}
                className="text-white bg-blue-600 hover:bg-blue-700 focus:ring-4 focus:ring-blue-300 font-bold rounded-xl text-sm px-8 py-3 transition-colors flex items-center gap-2 shadow-lg shadow-blue-200"
              >
                Tovább <ChevronRight className="w-5 h-5" />
              </button>
            </div>
          </>
        )}

        {activeTab === "attributes" && (
          <>
            {getActiveFiltersCount() > 0 && (
              <div className="mb-6 p-4 bg-blue-50 rounded-xl border border-blue-100">
                <p className="text-sm text-blue-800 font-medium mb-2">Kiválasztott kapcsolatok:</p>
                <div className="flex flex-wrap gap-2">
                  {filters.category_id && (
                    <span className="bg-white text-blue-700 text-xs px-3 py-1 rounded-full border border-blue-200">
                      {categories.find(c => Number(c.id) === filters.category_id)?.name}
                    </span>
                  )}
                  {filters.sub_category_id && (
                    <span className="bg-white text-blue-700 text-xs px-3 py-1 rounded-full border border-blue-200">
                      {subCategories.find(s => Number(s.id) === filters.sub_category_id)?.name}
                    </span>
                  )}
                  {filters.type_id && (
                    <span className="bg-white text-blue-700 text-xs px-3 py-1 rounded-full border border-blue-200">
                      {productTypes.find(t => Number(t.id) === filters.type_id)?.name}
                    </span>
                  )}
                  {filters.storing_condition_id && (
                    <span className="bg-white text-blue-700 text-xs px-3 py-1 rounded-full border border-blue-200">
                      {storingConditions.find(sc => Number(sc.id) === filters.storing_condition_id)?.description}
                    </span>
                  )}
                  {filters.brand_id && (
                    <span className="bg-white text-blue-700 text-xs px-3 py-1 rounded-full border border-blue-200">
                      {brands.find(b => Number(b.id) === filters.brand_id)?.name}
                    </span>
                  )}
                </div>
              </div>
            )}

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label className="text-xs font-bold text-slate-700 uppercase tracking-wide block mb-2">
                  Raktáron lévő mennyiség (Min - Max)
                </label>
                <div className="flex items-center gap-2">
                  <input
                    type="number"
                    className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg block p-3"
                    placeholder="Min"
                    value={filters.amount_min || ""}
                    onChange={(e) =>
                      setFilters((prev) => ({
                        ...prev,
                        amount_min: e.target.value ? Number(e.target.value) : undefined,
                      }))
                    }
                  />
                  <span className="text-gray-400">-</span>
                  <input
                    type="number"
                    className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg block p-3"
                    placeholder="Max"
                    value={filters.amount_max || ""}
                    onChange={(e) =>
                      setFilters((prev) => ({
                        ...prev,
                        amount_max: e.target.value ? Number(e.target.value) : undefined,
                      }))
                    }
                  />
                </div>
              </div>

              <div>
                <label className="text-xs font-bold text-slate-700 uppercase tracking-wide block mb-2">
                  Ár (Min - Max) Ft
                </label>
                <div className="flex items-center gap-2">
                  <input
                    type="number"
                    className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg block p-3"
                    placeholder="Min Ft"
                    value={filters.price_min || ""}
                    onChange={(e) =>
                      setFilters((prev) => ({
                        ...prev,
                        price_min: e.target.value ? Number(e.target.value) : undefined,
                      }))
                    }
                  />
                  <span className="text-gray-400">-</span>
                  <input
                    type="number"
                    className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg block p-3"
                    placeholder="Max Ft"
                    value={filters.price_max || ""}
                    onChange={(e) =>
                      setFilters((prev) => ({
                        ...prev,
                        price_max: e.target.value ? Number(e.target.value) : undefined,
                      }))
                    }
                  />
                </div>
              </div>

              <div>
                <label className="text-xs font-bold text-slate-700 uppercase tracking-wide block mb-2">
                  Kiszerelés méret
                </label>
                <input
                  type="number"
                  step="0.01"
                  className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg block p-3"
                  placeholder="Pl. 0.5 vagy 2"
                  value={filters.size_val || ""}
                  onChange={(e) =>
                    setFilters((prev) => ({
                      ...prev,
                      size_val: e.target.value ? Number(e.target.value) : undefined,
                    }))
                  }
                />
              </div>

              <div>
                <label className="text-xs font-bold text-slate-700 uppercase tracking-wide block mb-2">
                  Méret típus
                </label>
                <select
                  className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg block p-3"
                  value={filters.size_type || ""}
                  onChange={(e) =>
                    setFilters((prev) => ({ ...prev, size_type: e.target.value }))
                  }
                >
                  <option value="">Mindegy</option>
                  <option value="L">L</option>
                  <option value="kg">kg</option>
                  <option value="g">g</option>
                  <option value="db">db</option>
                </select>
              </div>

              <div className="md:col-span-2">
                <label className="text-xs font-bold text-slate-700 uppercase tracking-wide block mb-2">
                  Egyéb opciók
                </label>
                <div className="flex flex-wrap gap-4 p-4 bg-gray-50 rounded-lg border border-gray-200">
                  <label className="flex items-center cursor-pointer">
                    <input
                      type="checkbox"
                      className="w-5 h-5 text-blue-600 bg-white border-gray-300 rounded focus:ring-blue-500"
                      checked={filters.show_expired}
                      onChange={(e) =>
                        setFilters((prev) => ({ ...prev, show_expired: e.target.checked }))
                      }
                    />
                    <span className="ml-2 text-sm text-gray-700">
                      Csak lejárt termékek
                    </span>
                  </label>
                  <label className="flex items-center cursor-pointer">
                    <input
                      type="checkbox"
                      className="w-5 h-5 text-blue-600 bg-white border-gray-300 rounded focus:ring-blue-500"
                      checked={filters.has_warranty}
                      onChange={(e) =>
                        setFilters((prev) => ({ ...prev, has_warranty: e.target.checked }))
                      }
                    />
                    <span className="ml-2 text-sm text-gray-700">
                      Csak garanciális
                    </span>
                  </label>
                  <label className="flex items-center cursor-pointer">
                    <input
                      type="checkbox"
                      className="w-5 h-5 text-blue-600 bg-white border-gray-300 rounded focus:ring-blue-500"
                      checked={filters.is_discounted}
                      onChange={(e) =>
                        setFilters((prev) => ({ ...prev, is_discounted: e.target.checked }))
                      }
                    />
                    <span className="ml-2 text-sm text-gray-700">
                      Csak akciós
                    </span>
                  </label>
                </div>
              </div>

              <div className="md:col-span-2">
                <label className="text-xs font-bold text-slate-700 uppercase tracking-wide block mb-2">
                  Leírás keresés
                </label>
                <input
                  type="text"
                  className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-3"
                  value={filters.description || ""}
                  onChange={(e) =>
                    setFilters((prev) => ({ ...prev, description: e.target.value }))
                  }
                  placeholder="Keresés a termék leírásában..."
                />
              </div>
            </div>

            <div className="mt-8 pt-6 border-t border-gray-100 flex justify-between">
              <button
                type="button"
                onClick={goToPrevTab}
                className="text-gray-600 bg-gray-100 hover:bg-gray-200 focus:ring-4 focus:ring-gray-300 font-bold rounded-xl text-sm px-6 py-3 transition-colors flex items-center gap-2"
              >
                <ChevronLeft className="w-5 h-5" /> Vissza
              </button>
              <button
                type="submit"
                className="text-white bg-blue-600 hover:bg-blue-700 focus:ring-4 focus:ring-blue-300 font-bold rounded-xl text-sm px-8 py-3 transition-colors flex items-center gap-2 shadow-lg shadow-blue-200"
                disabled={!!loading}
              >
                {loading === "local" ? (
                  "Keresés..."
                ) : (
                  <>
                    <Filter className="w-5 h-5" /> Keresés indítása
                  </>
                )}
              </button>
              <button
                type="button"
                className="text-blue-600 border-2 border-blue-600 hover:bg-blue-50 focus:ring-4 focus:ring-blue-300 font-bold rounded-xl text-sm px-8 py-3 transition-colors flex items-center gap-2"
                disabled={!!loading}
                onClick={handleNetworkSearch}
              >
                {loading === "network" ? (
                  "Keresés..."
                ) : (
                  <>
                    <ScanSearch className="w-5 h-5" /> Hálózatos keresés indítása
                  </>
                )}
              </button>
            </div>
          </>
        )}
      </form>
    </div>
  );
}
