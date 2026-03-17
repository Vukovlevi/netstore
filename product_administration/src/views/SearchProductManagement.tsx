import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import SearchProductForm from "../components/forms/SearchProductForm";
import { searchService } from "../services/searchFilters";
import type { SearchFilters } from "../services/searchFilters";
import { categoryService } from "../services/categoryService";
import { subCategoryService } from "../services/subCategoryService";
import { productTypeService } from "../services/productTypeService";
import { brandService } from "../services/brandService";
import { storingConditionService } from "../services/storingConditionService";
import type {
  Category,
  SubCategory,
  ProductType,
  Brand,
  Product,
  StoringCondition,
  NetworkStoreResult,
} from "../types/Types";
import { ChevronLeft, ChevronRight, Store } from "lucide-react";
import FeedbackMessage from "../components/ui/FeedbackMessage";

export default function SearchProductManagement() {
  const navigate = useNavigate();
  const [results, setResults] = useState<Product[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState<"local" | "network" | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [searched, setSearched] = useState(false);
  const [activeTab, setActiveTab] = useState<"relations" | "attributes">(
    "relations",
  );

  const [categories, setCategories] = useState<Category[]>([]);
  const [subCategories, setSubCategories] = useState<SubCategory[]>([]);
  const [productTypes, setProductTypes] = useState<ProductType[]>([]);
  const [brands, setBrands] = useState<Brand[]>([]);
  const [storingConditions, setStoringConditions] = useState<
    StoringCondition[]
  >([]);

  const [filters, setFilters] = useState<SearchFilters>({
    page: 1,
    is_discounted: false,
    show_expired: false,
    has_warranty: false,
  });

  const [networkResults, setNetworkResults] = useState<NetworkStoreResult[]>([]);
  const [selectedStoreIndex, setSelectedStoreIndex] = useState<number>(-1);
  const [isNetworkMode, setIsNetworkMode] = useState(false);

  useEffect(() => {
    loadDependencies();
  }, []);

  const loadDependencies = async () => {
    try {
      const [cats, subs, types, brs, conds] = await Promise.all([
        categoryService.getAll().catch(() => []),
        subCategoryService.getAll().catch(() => []),
        productTypeService.getAll().catch(() => []),
        brandService.getAll().catch(() => []),
        storingConditionService.getAll().catch(() => []),
      ]);
      setCategories(Array.isArray(cats) ? cats : []);
      setSubCategories(Array.isArray(subs) ? subs : []);
      setProductTypes(Array.isArray(types) ? types : []);
      setBrands(Array.isArray(brs) ? brs : []);
      setStoringConditions(Array.isArray(conds) ? conds : []);
    } catch (err) {
      setError("Hiba történt az adatok betöltésekor.");
    }
  };

  const handleSearch = async (e?: React.FormEvent, pageOverride?: number) => {
    if (e) e.preventDefault();
    setLoading("local");
    setError(null);
    setSearched(true);
    setIsNetworkMode(false);
    setNetworkResults([]);
    setSelectedStoreIndex(-1);

    const currentFilters = { ...filters, page: pageOverride || filters.page };

    try {
      const response = await searchService.search(currentFilters);
      setResults(response.data);
      setTotal(response.total);
      setFilters((prev) => ({ ...prev, page: response.page }));
    } catch (err) {
      setError("Hiba történt a keresés során.");
    } finally {
      setLoading(null);
    }
  };

  const handleNetworkSearch = async (e?: React.FormEvent) => {
    if (e) e.preventDefault();
    setLoading("network");
    setError(null);
    setSearched(true);
    setIsNetworkMode(true);
    setResults([]);
    setTotal(0);
    setSelectedStoreIndex(-1);

    const currentFilters = { ...filters, page: 1 };

    try {
      const response = await searchService.networkSearch(currentFilters);
      setNetworkResults(response);
      if (response.length === 0) {
        setError("Nem érkezett válasz más üzletektől.");
      }
    } catch (err) {
      setError("Hiba történt a hálózatos keresés során.");
    } finally {
      setLoading(null);
    }
  };

  const handlePageChange = (newPage: number) => {
    setFilters((prev) => ({ ...prev, page: newPage }));
    handleSearch(undefined, newPage);
  };

  const clearFilters = () => {
    setFilters({
      page: 1,
      is_discounted: false,
      show_expired: false,
      has_warranty: false,
      name: "",
      category_id: undefined,
      sub_category_id: undefined,
      type_id: undefined,
      brand_id: undefined,
      storing_condition_id: undefined,
      amount_min: undefined,
      amount_max: undefined,
      size_val: undefined,
      size_type: "",
      price_min: undefined,
      price_max: undefined,
      description: "",
    });
    setResults([]);
    setTotal(0);
    setSearched(false);
    setIsNetworkMode(false);
    setNetworkResults([]);
    setSelectedStoreIndex(-1);
  };

  const totalPages = Math.ceil(total / 25);

  const selectedStore = selectedStoreIndex >= 0 ? networkResults[selectedStoreIndex] : null;
  const networkProducts: Product[] = selectedStore?.products?.data || [];
  const networkTotal: number = selectedStore?.products?.total || 0;

  return (
    <div className="flex flex-col gap-6 w-full max-w-6xl mx-auto">
      <SearchProductForm
        categories={categories}
        subCategories={subCategories}
        productTypes={productTypes}
        brands={brands}
        storingConditions={storingConditions}
        filters={filters}
        setFilters={setFilters}
        activeTab={activeTab}
        setActiveTab={setActiveTab}
        handleSearch={handleSearch}
        handleNetworkSearch={handleNetworkSearch}
        clearFilters={clearFilters}
        loading={loading}
      />

      {error && <FeedbackMessage type="error" message={error} />}

      {isNetworkMode && networkResults.length > 0 && (
        <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
          <div className="p-4 border-b border-gray-100 bg-gray-50 flex items-center gap-2">
            <Store className="w-5 h-5 text-blue-600" />
            <h2 className="font-bold text-slate-800">
              Üzlet kiválasztása ({networkResults.length} üzlet válaszolt)
            </h2>
          </div>
          <div className="p-4">
            <select
              className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-3"
              value={selectedStoreIndex}
              onChange={(e) => setSelectedStoreIndex(Number(e.target.value))}
            >
              <option value={-1}>Válasszon üzletet...</option>
              {networkResults.map((store, index) => (
                <option key={index} value={index}>
                  {store.store_detail.address}
                  {store.store_detail.storeTypeName
                    ? ` (${store.store_detail.storeTypeName})`
                    : ""}
                  {" "}- {store.products?.total || 0} termék
                </option>
              ))}
            </select>
          </div>

          {/* Store open hours */}
          {selectedStore && selectedStore.open_hours && selectedStore.open_hours.length > 0 && (
            <div className="px-4 pb-4">
              <p className="text-xs font-bold text-slate-700 uppercase tracking-wide mb-2">
                Nyitvatartás
              </p>
              <div className="flex flex-wrap gap-2">
                {selectedStore.open_hours.map((oh, i) => (
                  <span
                    key={i}
                    className="bg-blue-50 text-blue-700 text-xs px-3 py-1 rounded-full border border-blue-200"
                  >
                    {oh.day_name}: {oh.open} - {oh.close}
                  </span>
                ))}
              </div>
            </div>
          )}
        </div>
      )}

      {/* Network search: product list for selected store (read-only) */}
      {isNetworkMode && selectedStore && (
        <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden mb-10">
          <div className="p-4 border-b border-gray-100 flex justify-between items-center bg-gray-50">
            <h2 className="font-bold text-slate-800">
              Találatok - {selectedStore.store_detail.address} ({networkTotal} termék)
            </h2>
            <span className="text-xs text-gray-500">
              {networkTotal > 0 ? `${networkTotal} termék` : "0 / 0"}
            </span>
          </div>

          {networkProducts.length > 0 ? (
            <div className="overflow-x-auto">
              <table className="w-full text-sm text-left text-gray-500">
                <thead className="text-xs text-gray-700 uppercase bg-gray-50">
                  <tr>
                    <th className="px-3 py-2 md:px-6 md:py-3">Név</th>
                    <th className="px-3 py-2 md:px-6 md:py-3">Márka</th>
                    <th className="hidden md:table-cell px-6 py-3">Kategória</th>
                    <th className="hidden md:table-cell px-6 py-3">Típus</th>
                    <th className="px-3 py-2 md:px-6 md:py-3 text-right">Ár</th>
                    <th className="hidden md:table-cell px-6 py-3 text-right">Mennyiség</th>
                  </tr>
                </thead>
                <tbody>
                  {networkProducts.map((product, idx) => (
                    <tr
                      key={idx}
                      className="bg-white border-b transition-colors"
                    >
                      <td className="px-3 py-3 md:px-6 md:py-4 font-medium text-gray-900">
                        {product.name}
                      </td>
                      <td className="px-3 py-3 md:px-6 md:py-4">{product.brand_name}</td>
                      <td className="hidden md:table-cell px-6 py-4">
                        {product.category_name} / {product.sub_category_name}
                      </td>
                      <td className="hidden md:table-cell px-6 py-4">{product.type_name}</td>
                      <td className="px-3 py-3 md:px-6 md:py-4 text-right font-bold text-slate-700">
                        {product.discount > 0 && (
                          <span className="ml-2 bg-red-100 text-red-800 text-xs font-medium px-2 py-0.5 rounded mx-1">
                            -{Math.round(product.discount * 100)}%
                          </span>
                        )}
                        {product.price} Ft
                      </td>
                      <td className="hidden md:table-cell px-6 py-4 text-right">
                        {product.amount} db ({product.size} {product.size_type})
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          ) : (
            <div className="p-12 text-center text-gray-500">
              Nem található termék a megadott feltételekkel ebben az üzletben.
            </div>
          )}
        </div>
      )}

      {/* Local search results */}
      {!isNetworkMode && searched && (
        <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden mb-10">
          <div className="p-4 border-b border-gray-100 flex justify-between items-center bg-gray-50">
            <h2 className="font-bold text-slate-800">
              Találatok ({total} termék)
            </h2>
            <span className="text-xs text-gray-500">
              {total > 0
                ? `${(filters.page - 1) * 25 + 1} - ${Math.min(filters.page * 25, total)} / ${total}`
                : "0 / 0"}
            </span>
          </div>

          {results.length > 0 ? (
            <>
              <div className="overflow-x-auto">
                <table className="w-full text-sm text-left text-gray-500">
                  <thead className="text-xs text-gray-700 uppercase bg-gray-50">
                    <tr>
                      <th className="px-3 py-2 md:px-6 md:py-3">Név</th>
                      <th className="px-3 py-2 md:px-6 md:py-3">Márka</th>
                      <th className="hidden md:table-cell px-6 py-3">Kategória</th>
                      <th className="hidden md:table-cell px-6 py-3">Típus</th>
                      <th className="px-3 py-2 md:px-6 md:py-3 text-right">Ár</th>
                      <th className="hidden md:table-cell px-6 py-3 text-right">Mennyiség</th>
                    </tr>
                  </thead>
                  <tbody>
                    {results.map((product) => (
                      <tr
                        key={product.id}
                        onClick={() => {
                          sessionStorage.setItem("selectProductId", String(product.id));
                          navigate("/products");
                        }}
                        className="bg-white border-b hover:bg-blue-50 cursor-pointer transition-colors"
                      >
                        <td className="px-3 py-3 md:px-6 md:py-4 font-medium text-gray-900">
                          {product.name}
                        </td>
                        <td className="px-3 py-3 md:px-6 md:py-4">{product.brand_name}</td>
                        <td className="hidden md:table-cell px-6 py-4">
                          {product.category_name} / {product.sub_category_name}
                        </td>
                        <td className="hidden md:table-cell px-6 py-4">{product.type_name}</td>
                        <td className="px-3 py-3 md:px-6 md:py-4 text-right font-bold text-slate-700">
                          {product.discount > 0 && (
                            <span className="ml-2 bg-red-100 text-red-800 text-xs font-medium px-2 py-0.5 rounded mx-1">
                              -{Math.round(product.discount * 100)}%
                            </span>
                          )}
                          {product.price} Ft
                        </td>
                        <td className="hidden md:table-cell px-6 py-4 text-right">
                          {product.amount} db ({product.size} {product.size_type})
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              {totalPages > 1 && (
                <div className="p-4 border-t border-gray-100 flex items-center justify-center gap-2">
                  <button
                    onClick={() => handlePageChange(filters.page - 1)}
                    disabled={filters.page === 1}
                    className="p-2 rounded-lg hover:bg-gray-100 disabled:opacity-30 disabled:hover:bg-transparent"
                  >
                    <ChevronLeft className="w-5 h-5" />
                  </button>

                  {Array.from({ length: totalPages }, (_, i) => i + 1).map(
                    (pageNum) => (
                      <button
                        key={pageNum}
                        onClick={() => handlePageChange(pageNum)}
                        className={`w-8 h-8 rounded-lg text-sm font-medium transition-colors ${
                          pageNum === filters.page
                            ? "bg-blue-600 text-white shadow-md shadow-blue-200"
                            : "text-gray-600 hover:bg-gray-100"
                        }`}
                      >
                        {pageNum}
                      </button>
                    ),
                  )}

                  <button
                    onClick={() => handlePageChange(filters.page + 1)}
                    disabled={filters.page === totalPages}
                    className="p-2 rounded-lg hover:bg-gray-100 disabled:opacity-30 disabled:hover:bg-transparent"
                  >
                    <ChevronRight className="w-5 h-5" />
                  </button>
                </div>
              )}
            </>
          ) : (
            <div className="p-12 text-center text-gray-500">
              Nem található termék a megadott feltételekkel.
            </div>
          )}
        </div>
      )}
    </div>
  );
}
