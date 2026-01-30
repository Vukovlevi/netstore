import React, { useState, useEffect } from "react";
import ProductForm from "../components/forms/ProductForm";
import { productService } from "../services/productService";
import { categoryService } from "../services/categoryService";
import { subCategoryService } from "../services/subCategoryService";
import { productTypeService } from "../services/productTypeService";
import { brandService } from "../services/brandService";
import type {
  Category,
  SubCategory,
  ProductType,
  Brand,
  Product,
} from "../types/Types";
import { Search, PackageOpen } from "lucide-react";

export default function ProductManagement() {
  const [products, setProducts] = useState<Product[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [subCategories, setSubCategories] = useState<SubCategory[]>([]);
  const [productTypes, setProductTypes] = useState<ProductType[]>([]);
  const [brands, setBrands] = useState<Brand[]>([]);

  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [amount, setAmount] = useState<number | "">("");
  const [size, setSize] = useState("");
  const [sizeType, setSizeType] = useState("");
  const [expiresAt, setExpiresAt] = useState("");
  const [price, setPrice] = useState<number | "">("");
  const [discount, setDiscount] = useState<number | "">("");
  const [warranty, setWarranty] = useState("");
  const [categoryId, setCategoryId] = useState<number | "">("");
  const [subCategoryId, setSubCategoryId] = useState<number | "">("");
  const [typeId, setTypeId] = useState<number | "">("");
  const [brandId, setBrandId] = useState<number | "">("");

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);
  const [searchTerm, setSearchTerm] = useState("");

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    let prods: Product[] = [];
    try {
      const result = await productService.getAll();
      prods = Array.isArray(result) ? result : [];
      setProducts(prods);
    } catch (err) {
      console.error("Failed to load products");
      setProducts([]);
    }

    let subs: SubCategory[] = [];
    let types: ProductType[] = [];
    try {
      const [cats, subsResult, typesResult, brs] = await Promise.all([
        categoryService.getAll(),
        subCategoryService.getAll(),
        productTypeService.getAll(),
        brandService.getAll(),
      ]);
      subs = Array.isArray(subsResult) ? subsResult : [];
      types = Array.isArray(typesResult) ? typesResult : [];
      setCategories(Array.isArray(cats) ? cats : []);
      setSubCategories(subs);
      setProductTypes(types);
      setBrands(Array.isArray(brs) ? brs : []);
    } catch (err) {
      console.error("Failed to load dependencies");
    }

    const selectId = sessionStorage.getItem("selectProductId");
    if (selectId) {
      sessionStorage.removeItem("selectProductId");
      const product = prods.find((p) => Number(p.id) === Number(selectId));
      if (product) {
        setSelectedId(product.id);
        setName(product.name);
        setDescription(product.description);
        setAmount(product.amount);
        setSize(product.size);
        setSizeType(product.size_type);
        setExpiresAt(formatDate(product.expires_at));
        setPrice(product.price);
        setDiscount(product.discount);
        setWarranty(formatDate(product.warranty));
        setBrandId(product.brand_id);
        const type = types.find((t) => Number(t.id) === Number(product.type_id));
        if (type) {
          setTypeId(type.id);
          const sub = subs.find((s) => Number(s.id) === Number(type.sub_id));
          if (sub) {
            setSubCategoryId(sub.id);
            setCategoryId(sub.category_id);
          }
        }
      }
    }
  };

  const formatDate = (dateString: any): string => {
    if (!dateString) return "";
    const str = String(dateString);
    if (/^\d{4}-\d{2}-\d{2}$/.test(str)) return str;
    if (str.includes("T")) return str.split("T")[0];
    if (/^\d{4}$/.test(str)) return `${str}-01-01`;
    return "";
  };

  const selectProduct = (id: number) => {
    const product = products.find((p) => p.id === id);
    if (!product) return;

    setSelectedId(product.id);
    setName(product.name);
    setDescription(product.description);
    setAmount(product.amount);
    setSize(product.size);
    setSizeType(product.size_type);
    setExpiresAt(formatDate(product.expires_at));
    setPrice(product.price);
    setDiscount(product.discount);
    setWarranty(formatDate(product.warranty));
    setBrandId(product.brand_id);

    const type = productTypes.find((t) => t.id === product.type_id);
    if (type) {
      setTypeId(type.id);
      const sub = subCategories.find((s) => s.id === type.sub_id);
      if (sub) {
        setSubCategoryId(sub.id);
        setCategoryId(sub.category_id);
      } else {
        setSubCategoryId("");
        setCategoryId("");
      }
    } else {
      setTypeId("");
      setSubCategoryId("");
      setCategoryId("");
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (
      !typeId ||
      !brandId ||
      amount === "" ||
      !size.trim() ||
      !sizeType.trim() ||
      price === "" ||
      discount === ""
    )
      return;

    setError(null);
    setSuccessMsg(null);
    setLoading(true);

    const productData = {
      name,
      description,
      amount: Number(amount),
      size,
      size_type: sizeType,
      expires_at: formatDate(expiresAt) || null,
      price: Number(price),
      discount: Number(discount),
      warranty: formatDate(warranty) || null,
      type_id: Number(typeId),
      brand_id: Number(brandId),
    };

    try {
      if (selectedId) {
        await productService.update(selectedId, productData);
        setSuccessMsg("Termék sikeresen frissítve!");
      } else {
        await productService.create(productData);
        setSuccessMsg("Új termék létrehozva!");
        setName("");
        setDescription("");
        setAmount("");
        setSize("");
        setSizeType("");
        setExpiresAt("");
        setPrice("");
        setDiscount("");
        setWarranty("");
        setCategoryId("");
        setSubCategoryId("");
        setTypeId("");
        setBrandId("");
      }
      loadData();
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!selectedId) return;
    if (!window.confirm("Biztosan törölni szeretné ezt a terméket?")) return;

    setLoading(true);
    try {
      await productService.delete(selectedId);
      setSuccessMsg("Termék törölve!");
      setName("");
      setDescription("");
      setAmount("");
      setSize("");
      setSizeType("");
      setExpiresAt("");
      setPrice("");
      setDiscount("");
      setWarranty("");
      setCategoryId("");
      setSubCategoryId("");
      setTypeId("");
      setBrandId("");
      setSelectedId(null);
      loadData();
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const filteredProducts = products.filter((product) => {
    if (!searchTerm) return true;
    const term = searchTerm.toLowerCase();
    const matchesName = product.name.toLowerCase().includes(term);
    const matchesBrand = product.brand_name?.toLowerCase().includes(term);
    const matchesType = product.type_name?.toLowerCase().includes(term);

    return matchesName || matchesBrand || matchesType;
  });

  const searchResults = searchTerm ? filteredProducts : [];

  return (
    <div className="flex flex-col items-center justify-start gap-6 w-full relative">
      <div className="w-full max-w-3xl relative z-20">
        <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
          <Search className="h-5 w-5 text-gray-400" />
        </div>
        <input
          type="text"
          className="w-full pl-11 pr-4 py-3 bg-white border border-gray-100 rounded-2xl shadow-sm text-sm focus:outline-none focus:ring-2 focus:ring-blue-100 focus:border-blue-500 transition-all placeholder-gray-400 text-gray-700"
          placeholder="Keresés név, márka vagy típus alapján..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />

        {searchTerm && (
          <div className="absolute top-full left-0 right-0 mt-2 bg-white rounded-xl shadow-xl border border-gray-100 overflow-hidden max-h-96 overflow-y-auto">
            {searchResults.length > 0 ? (
              searchResults.map((product) => (
                <div
                  key={product.id}
                  onClick={() => {
                    selectProduct(product.id);
                    setSearchTerm("");
                  }}
                  className="p-4 hover:bg-blue-50 cursor-pointer border-b border-gray-50 last:border-0 transition-colors group"
                >
                  <div className="flex justify-between items-center">
                    <div className="flex items-center gap-3">
                      <div className="p-2 bg-gray-50 rounded-lg group-hover:bg-white transition-colors">
                        <PackageOpen className="h-5 w-5 text-blue-500" />
                      </div>
                      <div>
                        <div className="font-semibold text-slate-800">
                          {product.name}
                        </div>
                        <div className="text-xs text-slate-500">
                          {product.type_name}
                        </div>
                      </div>
                    </div>
                    <span className="text-xs font-medium text-blue-600 bg-blue-50 px-2 py-1 rounded-md">
                      {product.brand_name}
                    </span>
                  </div>
                </div>
              ))
            ) : (
              <div className="p-8 text-center text-gray-500 text-sm">
                Nincs találat a keresett kifejezésre.
              </div>
            )}
          </div>
        )}
      </div>

      <ProductForm
        products={filteredProducts}
        categories={categories}
        subCategories={subCategories}
        productTypes={productTypes}
        brands={brands}
        selectedId={selectedId}
        name={name}
        description={description}
        amount={amount}
        size={size}
        sizeType={sizeType}
        expiresAt={expiresAt}
        price={price}
        discount={discount}
        warranty={warranty}
        categoryId={categoryId}
        subCategoryId={subCategoryId}
        typeId={typeId}
        brandId={brandId}
        loading={loading}
        error={error}
        successMsg={successMsg}
        setName={setName}
        setDescription={setDescription}
        setAmount={setAmount}
        setSize={setSize}
        setSizeType={setSizeType}
        setExpiresAt={setExpiresAt}
        setPrice={setPrice}
        setDiscount={setDiscount}
        setWarranty={setWarranty}
        setCategoryId={setCategoryId}
        setSubCategoryId={setSubCategoryId}
        setTypeId={setTypeId}
        setBrandId={setBrandId}
        setSelectedId={setSelectedId}
        handleSubmit={handleSubmit}
        handleDelete={handleDelete}
      />
    </div>
  );
}
