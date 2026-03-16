import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import type { DailyMeal } from '../types';
import { ChevronDown, ChevronUp, RefreshCw, ShoppingCart } from 'lucide-react';

interface DayCardProps {
  meal: DailyMeal;
  dayName: string;
  onRegenerate: () => void;
  onAddToShopping: () => void;
}

const dayColors = [
  'from-purple-500 to-pink-500',
  'from-blue-500 to-cyan-500',
  'from-green-500 to-emerald-500',
  'from-yellow-500 to-orange-500',
  'from-red-500 to-pink-500',
  'from-indigo-500 to-purple-500',
  'from-orange-500 to-red-500',
];

export const DayCard: React.FC<DayCardProps> = ({
  meal,
  dayName,
  onRegenerate,
  onAddToShopping,
}) => {
  const { t } = useTranslation();
  const [showRecipe, setShowRecipe] = useState(false);
  const colorClass = dayColors[meal.day_of_week];

  return (
    <div className="bg-white rounded-xl shadow-lg overflow-hidden hover:shadow-xl transition">
      <div className={`bg-gradient-to-r ${colorClass} p-4`}>
        <h3 className="text-white font-heading font-bold text-xl">{dayName}</h3>
      </div>

      <div className="p-6">
        <h4 className="text-2xl font-heading font-bold text-gray-900 mb-4">
          {meal.menu_name}
        </h4>

        <div className="mb-4">
          <button
            onClick={() => setShowRecipe(!showRecipe)}
            className="flex items-center space-x-2 text-primary hover:text-accent transition font-semibold"
          >
            {showRecipe ? <ChevronUp size={20} /> : <ChevronDown size={20} />}
            <span>{showRecipe ? t('meal.hideRecipe') : t('meal.showRecipe')}</span>
          </button>

          {showRecipe && (
            <div className="mt-3 p-4 bg-gray-50 rounded-lg">
              <p className="text-gray-700 whitespace-pre-wrap">{meal.recipe}</p>
            </div>
          )}
        </div>

        <div className="mb-6">
          <h5 className="font-semibold text-gray-900 mb-2">{t('meal.shoppingList')}</h5>
          <ul className="space-y-1">
            {meal.shopping_items?.map((item, index) => (
              <li key={index} className="text-sm text-gray-700">
                • {item.item_name} - {item.quantity}{item.unit && ` ${item.unit}`}
              </li>
            ))}
          </ul>
        </div>

        <div className="flex space-x-2">
          <button
            onClick={onRegenerate}
            className="flex-1 bg-primary text-white py-2 px-4 rounded-lg hover:bg-primary/90 transition flex items-center justify-center space-x-2 text-sm font-semibold"
          >
            <RefreshCw size={16} />
            <span>{t('meal.regenerate')}</span>
          </button>

          <button
            onClick={onAddToShopping}
            className="flex-1 bg-success text-white py-2 px-4 rounded-lg hover:bg-success/90 transition flex items-center justify-center space-x-2 text-sm font-semibold"
          >
            <ShoppingCart size={16} />
            <span>{t('meal.addToList')}</span>
          </button>
        </div>
      </div>
    </div>
  );
};
