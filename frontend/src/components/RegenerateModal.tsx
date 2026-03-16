import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { DailyMealOption } from '../types';
import { X, Check } from 'lucide-react';

interface RegenerateModalProps {
  isOpen: boolean;
  onClose: () => void;
  options: DailyMealOption[];
  loading: boolean;
  onSelect: (index: number) => void;
  dayName: string;
}

export const RegenerateModal: React.FC<RegenerateModalProps> = ({
  isOpen,
  onClose,
  options,
  loading,
  onSelect,
  dayName,
}) => {
  const { t } = useTranslation();
  const [selectedIndex, setSelectedIndex] = useState<number | null>(null);

  if (!isOpen) return null;

  const handleSelect = () => {
    if (selectedIndex !== null) {
      onSelect(selectedIndex);
    }
  };

  return (
    <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-2xl shadow-2xl max-w-6xl w-full max-h-[90vh] overflow-hidden">
        <div className="bg-gradient-to-r from-primary to-accent p-6 flex items-center justify-between">
          <h2 className="text-2xl font-heading font-bold text-white">
            {t('meal.chooseNewMeal', { day: dayName })}
          </h2>
          <button
            onClick={onClose}
            className="text-white hover:bg-white/20 p-2 rounded-lg transition"
          >
            <X size={24} />
          </button>
        </div>

        <div className="p-6 overflow-y-auto max-h-[calc(90vh-180px)]">
          {loading ? (
            <div className="text-center py-12">
              <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto mb-4"></div>
              <p className="text-gray-600">{t('meal.generatingOptions')}</p>
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
              {options.map((option, index) => (
                <div
                  key={index}
                  onClick={() => setSelectedIndex(index)}
                  className={`border-2 rounded-xl p-6 cursor-pointer transition ${
                    selectedIndex === index
                      ? 'border-primary bg-primary/5'
                      : 'border-gray-200 hover:border-primary/50'
                  }`}
                >
                  <div className="flex items-start justify-between mb-3">
                    <h3 className="text-lg font-heading font-bold text-gray-900 flex-1">
                      {option.menu_name}
                    </h3>
                    {selectedIndex === index && (
                      <div className="bg-primary text-white rounded-full p-1">
                        <Check size={16} />
                      </div>
                    )}
                  </div>

                  <div className="mb-4">
                    <p className="text-sm text-gray-700 line-clamp-6">
                      {option.recipe}
                    </p>
                  </div>

                  <div>
                    <h4 className="font-semibold text-gray-900 mb-2 text-sm">
                      {t('meal.shoppingList')}
                    </h4>
                    <ul className="space-y-1">
                      {option.shopping_items.slice(0, 5).map((item, idx) => (
                        <li key={idx} className="text-xs text-gray-600">
                          • {item.item_name}
                        </li>
                      ))}
                      {option.shopping_items.length > 5 && (
                        <li className="text-xs text-gray-500">
                          {t('meal.moreItems', { count: option.shopping_items.length - 5 })}
                        </li>
                      )}
                    </ul>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {!loading && (
          <div className="border-t p-6 flex justify-end space-x-4">
            <button
              onClick={onClose}
              className="px-6 py-3 bg-gray-200 text-gray-700 rounded-lg font-semibold hover:bg-gray-300 transition"
            >
              {t('meal.cancel')}
            </button>
            <button
              onClick={handleSelect}
              disabled={selectedIndex === null}
              className="px-6 py-3 bg-gradient-to-r from-primary to-accent text-white rounded-lg font-semibold hover:shadow-lg transform hover:scale-105 transition disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {t('meal.selectMeal')}
            </button>
          </div>
        )}
      </div>
    </div>
  );
};
