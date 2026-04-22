import React, { useState, useEffect } from 'react';
import CategoryForm from '../components/forms/CategoryForm';
import { categoryService } from '../services/categoryService';
import type { Category } from '../types/Types';
import { useAuth, ROLES } from '../context/AuthContext';
import AccessDenied from '../components/AccessDenied';

export default function CategoryManagement() {
  const { canWrite } = useAuth();
  const [categories, setCategories] = useState<Category[]>([]);
  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [name, setName] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);

  const hasAccess = canWrite('category');

  useEffect(() => {
    if (hasAccess) {
      loadCategories();
    }
  }, [hasAccess]);

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
        const deleted = await categoryService.checkDeleted(name);
        if (deleted) {
          const ok = window.confirm(
            `Létezik egy korábban törölt kategória ezen a néven. Szeretné visszaállítani?\n\nKattintson az OK-ra a visszaállításhoz, vagy a Mégse-re a megszakításhoz.`
          );
          if (ok) {
            await categoryService.restore(deleted.id, name);
            setSuccessMsg("Kategória visszaállítva!");
            setName('');
          } else {
            setError("Visszaállítás megszakítva.");
          }
        } else {
          try {
            await categoryService.create(name);
            setSuccessMsg("Új kategória létrehozva!");
            setName('');
          } catch (err: any) {
            if (err.restorable && err.restoreId) {
              const ok = window.confirm(
                `Létezik egy korábban törölt kategória ezen a néven. Szeretné visszaállítani?\n\nKattintson az OK-ra a visszaállításhoz, vagy a Mégse-re a megszakításhoz.`
              );
              if (ok) {
                await categoryService.restore(err.restoreId, name);
                setSuccessMsg("Kategória visszaállítva!");
                setName('');
              } else {
                setError("Visszaállítás megszakítva.");
              }
            } else {
              throw err;
            }
          }
        }
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

  if (!hasAccess) {
    return (
      <AccessDenied
        resource="category"
        requiredRoles={[ROLES.UZLETVEZETO, ROLES.RAKTARVEZETO]}
      />
    );
  }

  return (
    <div className="flex justify-center items-start">
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
    </div>
  );
}