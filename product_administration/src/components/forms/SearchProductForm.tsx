import React from "react";
import type {
  Category,
  SubCategory,
  ProductType,
  Brand,
  StoringCondition,
} from "../../types/Types";
import type { SearchFilters } from "../../services/searchFilters";
import { Filter, X, Layers, Tag } from "lucide-react";

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
  clearFilters: () => void;
  loading: boolean;
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
  clearFilters,
  loading,
}: SearchProductFormProps) {
  const handleCategoryChange = (catId: number) => {
    setFilters((prev) => ({
      ...prev,
      category_id: catId,
      sub_category_id: undefined,
      type_id: undefined,
    }));
  };

  const handleSubCategoryChange = (subId: number) => {
    const sub = subCategories.find((s) => s.id === subId);
    if (sub) {
      setFilters((prev) => ({
        ...prev,
        sub_category_id: subId,
        category_id: sub.category_id,
        type_id: undefined,
      }));
    } else {
      setFilters((prev) => ({
        ...prev,
        sub_category_id: undefined,
        type_id: undefined,
      }));
    }
  };

  const handleTypeChange = (typeId: number) => {
    const type = productTypes.find((t) => t.id === typeId);
    if (type) {
      const sub = subCategories.find((s) => s.id === type.sub_id);
      const subId = sub ? sub.id : undefined;
      const catId = sub ? sub.category_id : undefined;

      setFilters((prev) => ({
        ...prev,
        type_id: typeId,
        sub_category_id: subId,
        category_id: catId,
        storing_condition_id: type.storing_condition_id,
      }));
    } else {
      setFilters((prev) => ({ ...prev, type_id: undefined }));
    }
  };

  const filteredSubCategories = filters.category_id
    ? subCategories.filter(
        (sc) => Number(sc.category_id) === Number(filters.category_id),
      )
    : [];

  const filteredTypes = filters.sub_category_id
    ? productTypes.filter(
        (pt) => Number(pt.sub_id) === Number(filters.sub_category_id),
      )
    : [];

  return (
    <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
      <div className="p-6 border-b border-gray-100 flex items-center justify-between">
        <div className="flex items-center gap-2">
          <h1 className="text-2xl font-bold text-slate-900">Komplex szűrés</h1>
        </div>
        <button
          onClick={clearFilters}
          className="text-sm text-gray-500 hover:text-red-600 flex items-center gap-1 transition-colors"
        >
          <X className="w-4 h-4" /> Szűrők törlése
        </button>
      </div>

      <div className="flex border-b border-gray-100">
        <button
          className={`flex-1 py-4 text-sm font-medium text-center transition-colors flex items-center justify-center gap-2 ${
            activeTab === "relations"
              ? "text-blue-600 border-b-2 border-blue-600 bg-blue-50/50"
              : "text-gray-500 hover:text-gray-700 hover:bg-gray-50"
          }`}
          onClick={() => setActiveTab("relations")}
        >
          <Layers className="w-4 h-4" /> Kapcsolatok és Kategóriák
        </button>
        <button
          className={`flex-1 py-4 text-sm font-medium text-center transition-colors flex items-center justify-center gap-2 ${
            activeTab === "attributes"
              ? "text-blue-600 border-b-2 border-blue-600 bg-blue-50/50"
              : "text-gray-500 hover:text-gray-700 hover:bg-gray-50"
          }`}
          onClick={() => setActiveTab("attributes")}
        >
          <Tag className="w-4 h-4" /> Termék tulajdonságok
        </button>
      </div>

      <form onSubmit={handleSearch} className="p-6">
        {activeTab === "relations" && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label className="text-xs font-bold text-slate-700 uppercase tracking-wide block mb-2">
                Kategória
              </label>
              <select
                className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-3"
                value={filters.category_id || ""}
                onChange={(e) => handleCategoryChange(Number(e.target.value))}
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
                className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-3 disabled:opacity-50"
                value={filters.sub_category_id || ""}
                onChange={(e) =>
                  handleSubCategoryChange(Number(e.target.value))
                }
                disabled={!filters.category_id}
              >
                <option value="">Válasszon...</option>
                {filteredSubCategories.map((s) => (
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
                className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-3 disabled:opacity-50"
                value={filters.type_id || ""}
                onChange={(e) => handleTypeChange(Number(e.target.value))}
                disabled={!filters.sub_category_id}
              >
                <option value="">Válasszon...</option>
                {filteredTypes.map((t) => (
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
                className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-3 disabled:opacity-50"
                value={filters.storing_condition_id || ""}
                onChange={(e) =>
                  setFilters({
                    ...filters,
                    storing_condition_id: Number(e.target.value),
                  })
                }
                disabled={!!filters.type_id}
              >
                <option value="">Válasszon...</option>
                {storingConditions.map((sc) => (
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
                  setFilters({ ...filters, brand_id: Number(e.target.value) })
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
        )}

        {activeTab === "attributes" && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="md:col-span-2">
              <label className="text-xs font-bold text-slate-700 uppercase tracking-wide block mb-2">
                Szabad szöveges keresés (Név)
              </label>
              <input
                type="text"
                className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-3"
                value={filters.name || ""}
                onChange={(e) =>
                  setFilters({ ...filters, name: e.target.value })
                }
                placeholder="Pl. Tej, Kenyér..."
              />
            </div>

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
                    setFilters({
                      ...filters,
                      amount_min: Number(e.target.value),
                    })
                  }
                />
                <span className="text-gray-400">-</span>
                <input
                  type="number"
                  className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg block p-3"
                  placeholder="Max"
                  value={filters.amount_max || ""}
                  onChange={(e) =>
                    setFilters({
                      ...filters,
                      amount_max: Number(e.target.value),
                    })
                  }
                />
              </div>
            </div>

            <div>
              <label className="text-xs font-bold text-slate-700 uppercase tracking-wide block mb-2">
                Lejárt termékek mutatása
              </label>
              <div className="flex items-center h-[46px]">
                <input
                  type="checkbox"
                  className="w-5 h-5 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500"
                  checked={filters.show_expired}
                  onChange={(e) =>
                    setFilters({ ...filters, show_expired: e.target.checked })
                  }
                />
                <span className="ml-2 text-sm text-gray-700">
                  Csak a lejártak mutatása
                </span>
              </div>
            </div>

            <div>
              <label className="text-xs font-bold text-slate-700 uppercase tracking-wide block mb-2">
                Kiszerelés méret
              </label>
              <input
                type="number"
                step="0.1"
                className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg block p-3"
                placeholder="Pl. 0.5 vagy 2"
                value={filters.size_val || ""}
                onChange={(e) =>
                  setFilters({ ...filters, size_val: Number(e.target.value) })
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
                  setFilters({ ...filters, size_type: e.target.value })
                }
              >
                <option value="">Mindegy</option>
                <option value="L">L</option>
                <option value="kg">kg</option>
                <option value="g">g</option>
                <option value="XXL">XXL</option>
              </select>
            </div>

            <div>
              <label className="text-xs font-bold text-slate-700 uppercase tracking-wide block mb-2">
                Ár (Min - Max)
              </label>
              <div className="flex items-center gap-2">
                <input
                  type="number"
                  className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg block p-3"
                  placeholder="Min Ft"
                  value={filters.price_min || ""}
                  onChange={(e) =>
                    setFilters({
                      ...filters,
                      price_min: Number(e.target.value),
                    })
                  }
                />
                <span className="text-gray-400">-</span>
                <input
                  type="number"
                  className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg block p-3"
                  placeholder="Max Ft"
                  value={filters.price_max || ""}
                  onChange={(e) =>
                    setFilters({
                      ...filters,
                      price_max: Number(e.target.value),
                    })
                  }
                />
              </div>
            </div>

            <div>
              <label className="text-xs font-bold text-slate-700 uppercase tracking-wide block mb-2">
                Egyéb opciók
              </label>
              <div className="flex gap-6 h-[46px] items-center">
                <label className="flex items-center cursor-pointer">
                  <input
                    type="checkbox"
                    className="w-5 h-5 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500"
                    checked={filters.has_warranty}
                    onChange={(e) =>
                      setFilters({ ...filters, has_warranty: e.target.checked })
                    }
                  />
                  <span className="ml-2 text-sm text-gray-700">
                    Garanciális
                  </span>
                </label>
                <label className="flex items-center cursor-pointer">
                  <input
                    type="checkbox"
                    className="w-5 h-5 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500"
                    checked={filters.is_discounted}
                    onChange={(e) =>
                      setFilters({
                        ...filters,
                        is_discounted: e.target.checked,
                      })
                    }
                  />
                  <span className="ml-2 text-sm text-gray-700">Akciós</span>
                </label>
              </div>
            </div>

            <div className="md:col-span-2">
              <label className="text-xs font-bold text-slate-700 uppercase tracking-wide block mb-2">
                Egyéb tulajdonságok (szabad szöveges)
              </label>
              <input
                type="text"
                className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-3"
                value={filters.other_properties || ""}
                onChange={(e) =>
                  setFilters({ ...filters, other_properties: e.target.value })
                }
                placeholder="Pl. leírásban keresés..."
              />
            </div>
          </div>
        )}

        <div className="mt-8 pt-6 border-t border-gray-100 flex justify-end">
          <button
            type="submit"
            className="text-white bg-blue-600 hover:bg-blue-700 focus:ring-4 focus:ring-blue-300 font-bold rounded-xl text-sm px-8 py-3 transition-colors flex items-center gap-2 shadow-lg shadow-blue-200"
            disabled={loading}
          >
            {loading ? (
              "Keresés..."
            ) : (
              <>
                <Filter className="w-5 h-5" /> Keresés indítása
              </>
            )}
          </button>
        </div>
      </form>
    </div>
  );
}
