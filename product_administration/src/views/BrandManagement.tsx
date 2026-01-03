import React, { useState, useEffect } from 'react';
import BrandForm from '../components/forms/BrandForm';
import { brandService } from '../services/brandService';
import type { Brand } from '../types/Types';

export default function BrandManagement() {
  const [brands, setBrands] = useState<Brand[]>([]);
  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [name, setName] = useState('');
  const [isOwn, setIsOwn] = useState(false);
  const [isTemporary, setIsTemporary] = useState(false);
  
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      const data = await brandService.getAll();
      setBrands(data);
    } catch (err) {
      console.error("Failed to load brands");
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccessMsg(null);
    setLoading(true);

    try {
      if (selectedId) {
        await brandService.update(selectedId, name, isOwn, isTemporary);
        setSuccessMsg("Márka sikeresen frissítve!");
      } else {
        await brandService.create(name, isOwn, isTemporary);
        setSuccessMsg("Új márka létrehozva!");
        setName('');
        setIsOwn(false);
        setIsTemporary(false);
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
    if (!window.confirm("Biztosan törölni szeretné ezt a márkát?")) return;

    setLoading(true);
    try {
      await brandService.delete(selectedId);
      setSuccessMsg("Márka törölve!");
      setName('');
      setIsOwn(false);
      setIsTemporary(false);
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
        <BrandForm
            brands={brands}
            selectedId={selectedId}
            name={name}
            isOwn={isOwn}
            isTemporary={isTemporary}
            loading={loading}
            error={error}
            successMsg={successMsg}
            setName={setName}
            setIsOwn={setIsOwn}
            setIsTemporary={setIsTemporary}
            setSelectedId={setSelectedId}
            handleSubmit={handleSubmit}
            handleDelete={handleDelete}
        />
    </div>
  );
}