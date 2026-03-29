import React from "react";
import type {
  Category,
  SubCategory,
  ProductType,
  Brand,
  Product,
} from "../../types/Types";
import FeedbackMessage from "../ui/FeedbackMessage";

interface ProductFormProps {
  products: Product[];
  categories: Category[];
  subCategories: SubCategory[];
  productTypes: ProductType[];
  brands: Brand[];

  selectedId: number | null;
  name: string;
  description: string;
  amount: number | "";
  size: string;
  sizeType: string;
  expiresAt: string;
  price: number | "";
  discount: number | "";
  warranty: string;
  categoryId: number | "";
  subCategoryId: number | "";
  typeId: number | "";
  brandId: number | "";

  loading: boolean;
  error: string | null;
  successMsg: string | null;

  setName: (val: string) => void;
  setDescription: (val: string) => void;
  setAmount: (val: number | "") => void;
  setSize: (val: string) => void;
  setSizeType: (val: string) => void;
  setExpiresAt: (val: string) => void;
  setPrice: (val: number | "") => void;
  setDiscount: (val: number | "") => void;
  setWarranty: (val: string) => void;
  setCategoryId: (val: number | "") => void;
  setSubCategoryId: (val: number | "") => void;
  setTypeId: (val: number | "") => void;
  setBrandId: (val: number | "") => void;
  setSelectedId: (id: number | null) => void;

  handleSubmit: (e: React.FormEvent) => void;
  handleDelete: () => void;
}

export default function ProductForm({
  products,
  categories,
  subCategories,
  productTypes,
  brands,
  selectedId,
  name,
  description,
  amount,
  size,
  sizeType,
  expiresAt,
  price,
  discount,
  warranty,
  categoryId,
  subCategoryId,
  typeId,
  brandId,
  loading,
  error,
  successMsg,
  setName,
  setDescription,
  setAmount,
  setSize,
  setSizeType,
  setExpiresAt,
  setPrice,
  setDiscount,
  setWarranty,
  setCategoryId,
  setSubCategoryId,
  setTypeId,
  setBrandId,
  setSelectedId,
  handleSubmit,
  handleDelete,
}: ProductFormProps) {
  const handleReset = () => {
    setSelectedId(null);
    setName("");
    setDescription("");
    setAmount("");
    setSize("");
    setSizeType("");
    setExpiresAt("");
    setPrice("");
    setDiscount(0);
    setWarranty("");
    setCategoryId("");
    setSubCategoryId("");
    setTypeId("");
    setBrandId("");
  };

  const formatDate = (dateString: any): string => {
    if (!dateString) return "";
    const str = String(dateString);

    if (/^\d{4}-\d{2}-\d{2}$/.test(str)) return str;

    if (str.includes("T")) return str.split("T")[0];

    if (/^\d{4}$/.test(str)) return `${str}-01-01`;
    return "";
  };

  const handleProductSelect = (val: string) => {
    const id = Number(val);
    if (id === 0) {
      handleReset();
    } else {
      const prod = products.find((p) => Number(p.id) === id);
      if (prod) {
        setSelectedId(id);
        setName(prod.name);
        setDescription(prod.description);
        setAmount(Number(prod.amount));
        setSize(prod.size);
        setSizeType(prod.size_type);
        setExpiresAt(formatDate(prod.expires_at));
        setPrice(Number(prod.price));
        setDiscount(Number(prod.discount));
        setWarranty(formatDate(prod.warranty));
        setBrandId(Number(prod.brand_id));

        const type = productTypes.find(
          (t) => Number(t.id) === Number(prod.type_id)
        );
        if (type) {
          setTypeId(Number(type.id));
          const sub = subCategories.find(
            (s) => Number(s.id) === Number(type.sub_id)
          );
          if (sub) {
            setSubCategoryId(Number(sub.id));
            setCategoryId(Number(sub.category_id));
          }
        }
      }
    }
  };

  const filteredSubCategories = categoryId
    ? subCategories.filter(
        (sc) => Number(sc.category_id) === Number(categoryId)
      )
    : [];

  const filteredTypes = subCategoryId
    ? productTypes.filter((pt) => Number(pt.sub_id) === Number(subCategoryId))
    : [];

  return (
    <div className="bg-white rounded-2xl shadow-sm border border-gray-100 p-4 md:p-8 w-full max-w-3xl">
      <div className="mb-8 flex justify-between items-start">
        <div>
          <h1 className="text-2xl font-bold text-slate-900">
            {selectedId ? "Termék szerkesztése" : "Új termék"}
          </h1>
          <p className="text-gray-500 mt-1 text-sm">
            {selectedId
              ? "Szerkessze a kiválasztott terméket vagy törölje."
              : "Adjon hozzá új terméket."}
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
              onChange={(e) => handleProductSelect(e.target.value)}
              value={selectedId || 0}
            >
              <option value={0}>-- Új létrehozása --</option>
              {products.map((p) => (
                <option key={p.id} value={p.id}>
                  {p.name}
                </option>
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
                  selectedId
                    ? "bg-gray-200 cursor-not-allowed opacity-75"
                    : "bg-gray-50"
                }`}
                value={categoryId}
                onChange={(e) => {
                  setCategoryId(Number(e.target.value));
                  setSubCategoryId("");
                  setTypeId("");
                }}
                disabled={!!selectedId}
              >
                <option value="" disabled>
                  Válasszon...
                </option>
                {categories.map((cat) => (
                  <option key={cat.id} value={cat.id}>
                    {cat.name}
                  </option>
                ))}
              </select>
            </div>

            <div className="flex flex-col gap-2">
              <label className="text-xs font-bold text-slate-700 uppercase tracking-wide">
                Alkategória
              </label>
              <select
                className={`appearance-none w-full border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-3 outline-none transition-all ${
                  !categoryId || !!selectedId
                    ? "bg-gray-200 cursor-not-allowed opacity-75"
                    : "bg-gray-50"
                }`}
                value={subCategoryId}
                onChange={(e) => {
                  setSubCategoryId(Number(e.target.value));
                  setTypeId("");
                }}
                disabled={!categoryId || !!selectedId}
              >
                <option value="" disabled>
                  Válasszon...
                </option>
                {filteredSubCategories.map((sub) => (
                  <option key={sub.id} value={sub.id}>
                    {sub.name}
                  </option>
                ))}
              </select>
            </div>
          </div>

          <div className="flex flex-col gap-2">
            <label className="text-xs font-bold text-slate-700 uppercase tracking-wide">
              Terméktípus
            </label>
            <select
              className={`appearance-none w-full border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-3 outline-none transition-all ${
                !subCategoryId || !!selectedId
                  ? "bg-gray-200 cursor-not-allowed opacity-75"
                  : "bg-gray-50"
              }`}
              value={typeId}
              onChange={(e) => setTypeId(Number(e.target.value))}
              disabled={!subCategoryId || !!selectedId}
            >
              <option value="" disabled>
                Válasszon...
              </option>
              {filteredTypes.map((t) => (
                <option key={t.id} value={t.id}>
                  {t.name}
                </option>
              ))}
            </select>
          </div>

          <div className="flex flex-col gap-2">
            <label className="text-xs font-bold text-slate-700 uppercase tracking-wide">
              Márka
            </label>
            <select
              className="w-full bg-gray-50 border border-gray-200 text-gray-700 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-2.5 outline-none transition-all appearance-none"
              onChange={(e) => setBrandId(Number(e.target.value))}
              value={brandId || ""}
            >
              <option value="" disabled>
                Válasszon...
              </option>
              {brands.map((b) => (
                <option key={b.id} value={b.id}>
                  {b.name}
                </option>
              ))}
            </select>
          </div>

          <div className="flex flex-col gap-2">
            <label className="text-xs font-bold text-slate-700 uppercase tracking-wide">
              Termék neve
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
            <div className="flex flex-col gap-2">
              <label className="text-xs font-bold text-slate-700 uppercase tracking-wide">
                Lejárat
              </label>
              <input
                type="date"
                value={expiresAt}
                onChange={(e) => setExpiresAt(e.target.value)}
                className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-2 focus:ring-blue-100 focus:border-blue-500 block p-3 outline-none transition-all"
              />
            </div>
            <div className="flex flex-col gap-2">
              <label className="text-xs font-bold text-slate-700 uppercase tracking-wide">
                Ár (Ft)
              </label>
              <input
                type="number"
                value={price}
                onChange={(e) => setPrice(Number(e.target.value))}
                placeholder="pl. 450"
                className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-2 focus:ring-blue-100 focus:border-blue-500 block p-3 outline-none transition-all"
              />
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="flex flex-col gap-2">
              <label className="text-xs font-bold text-slate-700 uppercase tracking-wide">
                Mennyiség
              </label>
              <input
                type="number"
                value={amount}
                onChange={(e) => setAmount(Number(e.target.value))}
                placeholder="pl. 10"
                className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-2 focus:ring-blue-100 focus:border-blue-500 block p-3 outline-none transition-all"
              />
            </div>
            <div className="flex flex-col gap-2">
              <label className="text-xs font-bold text-slate-700 uppercase tracking-wide">
                Kiszerelés
              </label>
              <input
                type="text"
                value={size}
                onChange={(e) => setSize(e.target.value)}
                placeholder="pl. 0.5"
                className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-2 focus:ring-blue-100 focus:border-blue-500 block p-3 outline-none transition-all"
              />
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="flex flex-col gap-2">
              <label className="text-xs font-bold text-slate-700 uppercase tracking-wide">
                Kiszerelés Típusa
              </label>
              <input
                type="text"
                value={sizeType}
                onChange={(e) => setSizeType(e.target.value)}
                placeholder="pl. kg"
                className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-2 focus:ring-blue-100 focus:border-blue-500 block p-3 outline-none transition-all"
              />
            </div>
            <div className="flex flex-col gap-2">
              <label className="text-xs font-bold text-slate-700 uppercase tracking-wide">
                Kedvezmény (%)
              </label>
              <input
                type="number"
                value={discount}
                onChange={(e) => setDiscount(Number(e.target.value))}
                placeholder="pl. 10"
                className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-2 focus:ring-blue-100 focus:border-blue-500 block p-3 outline-none transition-all"
              />
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="flex flex-col gap-2">
              <label className="text-xs font-bold text-slate-700 uppercase tracking-wide">
                Garancia Lejárata
              </label>
              <input
                type="date"
                value={warranty}
                onChange={(e) => setWarranty(e.target.value)}
                className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-2 focus:ring-blue-100 focus:border-blue-500 block p-3 outline-none transition-all"
              />
            </div>
          </div>

          <div className="flex flex-col gap-2">
            <label className="text-xs font-bold text-slate-700 uppercase tracking-wide">
              Termék leírása
            </label>
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Rövid leírás a termékről..."
              rows={3}
              className="w-full bg-gray-50 border border-gray-200 text-gray-900 text-sm rounded-lg focus:ring-2 focus:ring-blue-100 focus:border-blue-500 block p-3 outline-none transition-all resize-none"
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
            disabled={
              loading ||
              !name.trim() ||
              !description.trim() ||
              amount === "" ||
              !size.trim() ||
              !sizeType.trim() ||
              price === "" ||
              discount === "" ||
              !typeId ||
              !brandId
            }
            className="w-full md:w-auto px-6 py-2.5 text-sm font-bold text-white bg-blue-600 rounded-lg hover:bg-blue-700 transition-all shadow-sm shadow-blue-200 focus:outline-none focus:ring-4 focus:ring-blue-100 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? "Mentés..." : selectedId ? "Módosítás" : "Létrehozás"}
          </button>
        </div>
      </form>
    </div>
  );
}
