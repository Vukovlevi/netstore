import React, { useState, useEffect } from 'react';
import SubCategoryForm from '../components/forms/SubCategoryForm';
import { subCategoryService } from '../services/subCategoryService';
import { categoryService } from '../services/categoryService';
import type { Category, SubCategory } from '../types/Types';

export default function SubCategoryManagement() {
  const [subCategories, setSubCategories] = useState<SubCategory[]>([]);
  const [categories, setCategories] = useState<Category[]>([]);
  
  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [name, setName] = useState('');
  const [categoryId, setCategoryId] = useState<number | ''>('');
  
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      const [subs, cats] = await Promise.all([
        subCategoryService.getAll(),
        categoryService.getAll()
      ]);
      setSubCategories(subs);
      setCategories(cats);
    } catch (err) {
      console.error("Failed to load data");
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!categoryId) return;
    
    setError(null);
    setSuccessMsg(null);
    setLoading(true);

    try {
      if (selectedId) {
        await subCategoryService.update(selectedId, name, Number(categoryId));
        setSuccessMsg("Alkategória sikeresen frissítve!");
      } else {
        await subCategoryService.create(name, Number(categoryId));
        setSuccessMsg("Új alkategória létrehozva!");
        setName('');
        setCategoryId('');
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
    if (!window.confirm("Biztosan törölni szeretné ezt az alkategóriát?")) return;

    setLoading(true);
    try {
      await subCategoryService.delete(selectedId);
      setSuccessMsg("Alkategória törölve!");
      setName('');
      setCategoryId('');
      setSelectedId(null);
      loadData();
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex justify-center items-start">
        <SubCategoryForm
            subCategories={subCategories}
            categories={categories}
            selectedId={selectedId}
            name={name}
            categoryId={categoryId}
            loading={loading}
            error={error}
            successMsg={successMsg}
            setName={setName}
            setCategoryId={setCategoryId}
            setSelectedId={setSelectedId}
            handleSubmit={handleSubmit}
            handleDelete={handleDelete}
        />
    </div>
  );
}