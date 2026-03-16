import React from 'react';
import type { ShoppingListItem as ShoppingItem } from '../types';
import { Trash2 } from 'lucide-react';

interface ShoppingListItemProps {
  item: ShoppingItem;
  onToggle: () => void;
  onDelete: () => void;
}

export const ShoppingListItem: React.FC<ShoppingListItemProps> = ({
  item,
  onToggle,
  onDelete,
}) => {

  return (
    <div
      className={`flex items-center gap-2 p-2 border-b transition ${
        item.is_checked ? 'bg-gray-50' : 'bg-white hover:bg-gray-50'
      }`}
    >
      <input
        type="checkbox"
        checked={item.is_checked}
        onChange={onToggle}
        className="w-4 h-4 text-primary rounded focus:ring-2 focus:ring-primary cursor-pointer flex-shrink-0"
      />

      <div className="flex-1 min-w-0 flex items-baseline gap-2">
        <p
          className={`font-medium text-sm ${
            item.is_checked ? 'text-gray-400 line-through' : 'text-gray-900'
          }`}
        >
          {item.item_name}
        </p>
        <span className={`text-xs whitespace-nowrap ${item.is_checked ? 'text-gray-300' : 'text-gray-500'}`}>
          {item.quantity}{item.unit && ` ${item.unit}`}
        </span>
      </div>

      <button
        onClick={onDelete}
        className="text-red-500 hover:text-red-700 hover:bg-red-50 p-1 rounded transition flex-shrink-0"
      >
        <Trash2 size={14} />
      </button>
    </div>
  );
};
