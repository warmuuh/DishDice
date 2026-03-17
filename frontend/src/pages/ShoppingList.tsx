import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { Header } from '../components/Header';
import { ShoppingListItem } from '../components/ShoppingListItem';
import { shoppingService } from '../services/shoppingService';
import type { ShoppingListItem as ShoppingItem } from '../types';
import { toast } from 'react-toastify';
import { Plus, Trash2, Eye, EyeOff } from 'lucide-react';

export const ShoppingList: React.FC = () => {
  const { t } = useTranslation();
  const [items, setItems] = useState<ShoppingItem[]>([]);
  const [loading, setLoading] = useState(true);
  const [showChecked, setShowChecked] = useState(false);
  const [itemName, setItemName] = useState('');
  const [quantity, setQuantity] = useState('');
  const [unit, setUnit] = useState('');
  const [addingItem, setAddingItem] = useState(false);

  useEffect(() => {
    loadItems();
  }, [showChecked]);

  const loadItems = async () => {
    try {
      const data = await shoppingService.getShoppingList(showChecked);
      setItems(data || []);
    } catch (error) {
      toast.error(t('shopping.loadFailed'));
    } finally {
      setLoading(false);
    }
  };

  const handleAddItem = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!itemName || !quantity) {
      toast.error(t('shopping.enterNameQuantity'));
      return;
    }

    setAddingItem(true);

    try {
      await shoppingService.addItem({ item_name: itemName, quantity, unit });
      toast.success(t('shopping.itemAdded'));
      setItemName('');
      setQuantity('');
      setUnit('');
      loadItems();
    } catch (error) {
      toast.error(t('shopping.addFailed'));
    } finally {
      setAddingItem(false);
    }
  };

  const handleToggleItem = async (id: string) => {
    try {
      await shoppingService.toggleItem(id);
      loadItems();
    } catch (error) {
      toast.error(t('shopping.updateFailed'));
    }
  };

  const handleDeleteItem = async (id: string) => {
    try {
      await shoppingService.deleteItem(id);
      toast.success(t('shopping.itemDeleted'));
      loadItems();
    } catch (error) {
      toast.error(t('shopping.deleteFailed'));
    }
  };

  const handleDeleteChecked = async () => {
    const checkedCount = items.filter((i) => i.is_checked).length;

    if (checkedCount === 0) {
      toast.error(t('shopping.noCheckedItems'));
      return;
    }

    if (!confirm(t('shopping.deleteConfirm', { count: checkedCount }))) {
      return;
    }

    try {
      await shoppingService.deleteChecked();
      toast.success(t('shopping.checkedDeleted'));
      loadItems();
    } catch (error) {
      toast.error(t('shopping.deleteFailed'));
    }
  };

  const handleDeleteAll = async () => {
    if (items.length === 0) {
      return;
    }

    if (!confirm(t('shopping.deleteAllConfirm', { count: items.length }))) {
      return;
    }

    try {
      // Delete all items one by one
      await Promise.all(items.map(item => shoppingService.deleteItem(item.id)));
      toast.success(t('shopping.allDeleted'));
      loadItems();
    } catch (error) {
      toast.error(t('shopping.deleteFailed'));
    }
  };

  const uncheckedCount = items.filter((i) => !i.is_checked).length;
  const checkedCount = items.filter((i) => i.is_checked).length;

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Header />
        <div className="flex items-center justify-center p-12">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary"></div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />

      <main className="container mx-auto px-3 py-4 md:px-4 md:py-6">
        <div className="max-w-2xl mx-auto">
          <div className="bg-white rounded-xl shadow-md p-4 md:p-6 mb-4">
            <div className="flex items-center justify-between mb-3">
              <h1 className="text-2xl md:text-3xl font-heading font-bold text-primary">
                {t('shopping.title')}
              </h1>
              <div className="text-xs md:text-sm text-gray-600">
                {t('shopping.status', { pending: uncheckedCount, completed: checkedCount })}
              </div>
            </div>

            <form onSubmit={handleAddItem} className="flex gap-2 mb-3">
              <input
                type="text"
                value={itemName}
                onChange={(e) => setItemName(e.target.value)}
                placeholder={t('shopping.itemName')}
                className="flex-1 px-3 py-2 text-sm border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
                disabled={addingItem}
              />
              <input
                type="text"
                value={quantity}
                onChange={(e) => setQuantity(e.target.value)}
                placeholder={t('shopping.quantity')}
                className="w-16 px-2 py-2 text-sm border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
                disabled={addingItem}
              />
              <input
                type="text"
                value={unit}
                onChange={(e) => setUnit(e.target.value)}
                placeholder={t('shopping.unit')}
                className="w-16 px-2 py-2 text-sm border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
                disabled={addingItem}
              />
              <button
                type="submit"
                disabled={addingItem}
                className="bg-gradient-to-r from-primary to-accent text-white px-3 md:px-4 py-2 rounded-lg font-semibold hover:shadow-lg transition disabled:opacity-50"
              >
                <Plus size={18} />
              </button>
            </form>

            <div className="flex flex-wrap items-center gap-2">
              <button
                onClick={() => setShowChecked(!showChecked)}
                className="flex items-center gap-1 text-xs text-primary hover:text-accent transition px-2 py-1 rounded hover:bg-primary/5"
              >
                {showChecked ? <EyeOff size={14} /> : <Eye size={14} />}
                <span className="hidden sm:inline">{showChecked ? t('shopping.hideChecked') : t('shopping.showChecked')}</span>
              </button>

              {checkedCount > 0 && (
                <button
                  onClick={handleDeleteChecked}
                  className="flex items-center gap-1 text-xs text-orange-600 hover:text-orange-800 transition px-2 py-1 rounded hover:bg-orange-50"
                >
                  <Trash2 size={14} />
                  <span className="hidden sm:inline">{t('shopping.clearChecked')}</span>
                </button>
              )}

              {items.length > 0 && (
                <button
                  onClick={handleDeleteAll}
                  className="flex items-center gap-1 text-xs text-red-600 hover:text-red-800 transition px-2 py-1 rounded hover:bg-red-50 ml-auto"
                >
                  <Trash2 size={14} />
                  <span>{t('shopping.deleteAll')}</span>
                </button>
              )}
            </div>
          </div>

          {items.length === 0 ? (
            <div className="bg-white rounded-xl shadow-md p-8 text-center">
              <div className="text-5xl mb-3">🛒</div>
              <h2 className="text-xl font-heading font-bold text-gray-900 mb-1">
                {t('shopping.empty')}
              </h2>
              <p className="text-sm text-gray-600">
                {t('shopping.emptyHint')}
              </p>
            </div>
          ) : (
            <div className="bg-white rounded-xl shadow-md overflow-hidden">
              {items.map((item) => (
                <ShoppingListItem
                  key={item.id}
                  item={item}
                  onToggle={() => handleToggleItem(item.id)}
                  onDelete={() => handleDeleteItem(item.id)}
                />
              ))}
            </div>
          )}
        </div>
      </main>
    </div>
  );
};
