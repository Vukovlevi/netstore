import React, { useState, useEffect } from 'react';
import StoringConditionForm from '../components/forms/StoringConditionForm';
import { storingConditionService } from '../services/storingConditionService';
import type { StoringCondition } from '../types/Types';

export default function StoringConditionManagement() {
  const [conditions, setConditions] = useState<StoringCondition[]>([]);
  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [description, setDescription] = useState('');
  
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [successMsg, setSuccessMsg] = useState<string | null>(null);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      const data = await storingConditionService.getAll();
      setConditions(data);
    } catch (err) {
      console.error("Failed to load storing conditions");
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccessMsg(null);
    setLoading(true);

    try {
      if (selectedId) {
        await storingConditionService.update(selectedId, description);
        setSuccessMsg("Tárolási körülmény sikeresen frissítve!");
      } else {
        await storingConditionService.create(description);
        setSuccessMsg("Új tárolási körülmény létrehozva!");
        setDescription('');
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
    if (!window.confirm("Biztosan törölni szeretné ezt a tárolási körülményt?")) return;

    setLoading(true);
    try {
      await storingConditionService.delete(selectedId);
      setSuccessMsg("Tárolási körülmény törölve!");
      setDescription('');
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
        <StoringConditionForm
            conditions={conditions}
            selectedId={selectedId}
            description={description}
            loading={loading}
            error={error}
            successMsg={successMsg}
            setDescription={setDescription}
            setSelectedId={setSelectedId}
            handleSubmit={handleSubmit}
            handleDelete={handleDelete}
        />
    </div>
  );
}