import React, { useState, useEffect } from 'react';
import ProductTypeForm from '../components/forms/ProductTypeForm';
import { productTypeService } from '../services/productTypeService';
import { categoryService } from '../services/categoryService';
import { subCategoryService } from '../services/subCategoryService';
import { storingConditionService } from '../services/storingConditionService';
import type { Category, SubCategory, ProductType, StoringCondition } from '../types/Types';
import { useAuth, ROLES } from '../context/AuthContext';
import AccessDenied from '../components/AccessDenied';

export default function ProductTypeManagement() {
  const { canWrite } = useAuth();
  const [productTypes, setProductTypes] = useState<ProductType[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  const [subCategories, setSubCategories] = useState<SubCategory[]>([]);
  const [storingConditions, setStoringConditions] = useState<StoringCondition[]>([]);

  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [categoryId, setCategoryId] = useState<number | ''>('');
  const [subCategoryId, setSubCategoryId] = useState<number | ''>('');
  const [storingConditionId, setStoringConditionId] = useState<number | ''>('');

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);

  const hasAccess = canWrite('product_type');

  useEffect(() => {
    if (hasAccess) {
      loadData();
    }
  }, [hasAccess]);

  const loadData = async () => {
    try {
      const [types, cats, subs, conds] = await Promise.all([
        productTypeService.getAll(),
        categoryService.getAll(),
        subCategoryService.getAll(),
        storingConditionService.getAll()
      ]);
      setProductTypes(types);
      setCategories(cats);
      setSubCategories(subs);
      setStoringConditions(conds);
    } catch (err) {
      console.error("Failed to load data");
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!subCategoryId || !storingConditionId) return;
    
    setError(null);
    setSuccessMsg(null);
    setLoading(true);

    try {
      if (selectedId) {
        await productTypeService.update(
            selectedId, 
            name, 
            description, 
            Number(subCategoryId),
            Number(storingConditionId)
        );
        setSuccessMsg("Terméktípus sikeresen frissítve!");
      } else {
        await productTypeService.create(
            name, 
            description, 
            Number(subCategoryId),
            Number(storingConditionId)
        );
        setSuccessMsg("Új terméktípus létrehozva!");
        setName('');
        setDescription('');
        setCategoryId('');
        setSubCategoryId('');
        setStoringConditionId('');
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
    if (!window.confirm("Biztosan törölni szeretné ezt a terméktípust?")) return;

    setLoading(true);
    try {
      await productTypeService.delete(selectedId);
      setSuccessMsg("Terméktípus törölve!");
      setName('');
      setDescription('');
      setCategoryId('');
      setSubCategoryId('');
      setStoringConditionId('');
      setSelectedId(null);
      loadData();
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  if (!hasAccess) {
    return (
      <AccessDenied
        resource="product_type"
        requiredRoles={[ROLES.UZLETVEZETO, ROLES.RAKTARVEZETO]}
      />
    );
  }

  return (
    <div className="flex justify-center items-start">
        <ProductTypeForm
            productTypes={productTypes}
            categories={categories}
            subCategories={subCategories}
            storingConditions={storingConditions}
            selectedId={selectedId}
            name={name}
            description={description}
            categoryId={categoryId}
            subCategoryId={subCategoryId}
            storingConditionId={storingConditionId}
            loading={loading}
            error={error}
            successMsg={successMsg}
            setName={setName}
            setDescription={setDescription}
            setCategoryId={setCategoryId}
            setSubCategoryId={setSubCategoryId}
            setStoringConditionId={setStoringConditionId}
            setSelectedId={setSelectedId}
            handleSubmit={handleSubmit}
            handleDelete={handleDelete}
        />
    </div>
  );
}