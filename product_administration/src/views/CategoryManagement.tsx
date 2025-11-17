import React, { useState, useEffect } from 'react';
import CategoryForm from '../components/forms/CategoryForm';
import { categoryService } from '../services/categoryService';
import type { Category } from '../types/Types';

export default function CategoryManagement() {
  const [categories, setCategories] = useState<Category[]>([]);
  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [name, setName] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);

  useEffect(() => {
    loadCategories();
  }, []);

  const loadCategories = async () => {
    try {
      const data = await categoryService.getAll();
      setCategories(data);
    } catch (err) {
      console.error("Failed to load categories");
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccessMsg(null);
    setLoading(true);

    try {
      if (selectedId) {
        await categoryService.update(selectedId, name);
        setSuccessMsg("Kategória sikeresen frissítve!");
      } else {
        await categoryService.create(name);
        setSuccessMsg("Új kategória létrehozva!");
        setName('');
      }
      loadCategories();
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!selectedId) return;
    if (!window.confirm("Biztosan törölni szeretné ezt a kategóriát?")) return;

    setLoading(true);
    try {
      await categoryService.delete(selectedId);
      setSuccessMsg("Kategória törölve!");
      setName('');
      setSelectedId(null);
      loadCategories();
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <main className="flex-1 ml-64 p-12 flex justify-center items-start">
        <CategoryForm
            categories={categories}
            selectedId={selectedId}
            name={name}
            loading={loading}
            error={error}
            successMsg={successMsg}
            setName={setName}
            setSelectedId={setSelectedId}
            handleSubmit={handleSubmit}
            handleDelete={handleDelete}
        />
    </main>
  );
}