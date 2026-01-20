import React, { useState, useEffect } from "react";
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
} from "../types/Types";
import { ChevronLeft, ChevronRight } from "lucide-react";
import FeedbackMessage from "../components/ui/FeedbackMessage";

export default function SearchProductManagement() {
  const [results, setResults] = useState<Product[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(false);
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

  useEffect(() => {
    loadDependencies();
  }, []);

  const loadDependencies = async () => {
    try {
      const [cats, subs, types, brs, conds] = await Promise.all([
        categoryService.getAll(),
        subCategoryService.getAll(),
        productTypeService.getAll(),
        brandService.getAll(),
        storingConditionService.getAll(),
      ]);
      setCategories(Array.isArray(cats) ? cats : []);
      setSubCategories(Array.isArray(subs) ? subs : []);
      setProductTypes(Array.isArray(types) ? types : []);
      setBrands(Array.isArray(brs) ? brs : []);
      setStoringConditions(Array.isArray(conds) ? conds : []);
    } catch (err) {
      console.error("Failed to load dependencies");
    }
  };

  const handleSearch = async (e?: React.FormEvent, pageOverride?: number) => {
    if (e) e.preventDefault();
    setLoading(true);
    setError(null);
    setSearched(true);

    const currentFilters = { ...filters, page: pageOverride || filters.page };

    try {
      const response = await searchService.search(currentFilters);
      setResults(response.data);
      setTotal(response.total);
      setFilters((prev) => ({ ...prev, page: response.page }));
    } catch (err) {
      setError("Hiba történt a keresés során.");
    } finally {
      setLoading(false);
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
      other_properties: "",
    });
    setResults([]);
    setTotal(0);
    setSearched(false);
  };

  const totalPages = Math.ceil(total / 25);

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
        clearFilters={clearFilters}
        loading={loading}
      />

      {error && <FeedbackMessage type="error" message={error} />}

      {searched && (
        <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden mb-10">
          <div className="p-4 border-b border-gray-100 flex justify-between items-center bg-gray-50">
            <h2 className="font-bold text-slate-800">
              Találatok ({total} termék)
            </h2>
            <span className="text-xs text-gray-500">
              {(filters.page - 1) * 25 + 1} -{" "}
              {Math.min(filters.page * 25, total)} / {total}
            </span>
          </div>

          {results.length > 0 ? (
            <>
              <div className="overflow-x-auto">
                <table className="w-full text-sm text-left text-gray-500">
                  <thead className="text-xs text-gray-700 uppercase bg-gray-50">
                    <tr>
                      <th className="px-6 py-3">Név</th>
                      <th className="px-6 py-3">Márka</th>
                      <th className="px-6 py-3">Kategória</th>
                      <th className="px-6 py-3">Típus</th>
                      <th className="px-6 py-3 text-right">Ár</th>
                      <th className="px-6 py-3 text-right">Mennyiség</th>
                    </tr>
                  </thead>
                  <tbody>
                    {results.map((product) => (
                      <tr
                        key={product.id}
                        className="bg-white border-b hover:bg-gray-50"
                      >
                        <td className="px-6 py-4 font-medium text-gray-900">
                          {product.name}
                        </td>
                        <td className="px-6 py-4">{product.brand_name}</td>
                        <td className="px-6 py-4">
                          {product.category_name} / {product.sub_category_name}
                        </td>
                        <td className="px-6 py-4">{product.type_name}</td>
                        <td className="px-6 py-4 text-right font-bold text-slate-700">
                          {product.price} Ft
                          {product.discount > 0 && (
                            <span className="ml-2 bg-red-100 text-red-800 text-xs font-medium px-2 py-0.5 rounded">
                              -{product.discount}%
                            </span>
                          )}
                        </td>
                        <td className="px-6 py-4 text-right">
                          {product.amount} {product.size}
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
