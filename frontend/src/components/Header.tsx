import React from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useAuth } from '../context/AuthContext';
import { LogOut, User, ShoppingCart, Home, Settings } from 'lucide-react';

export const Header: React.FC = () => {
  const { t } = useTranslation();
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/login');
  };

  return (
    <header className="bg-gradient-to-r from-primary to-accent text-white shadow-lg">
      <div className="container mx-auto px-4 py-4">
        <div className="flex items-center justify-between">
          <Link to="/dashboard" className="flex items-center space-x-2">
            <h1 className="text-2xl font-heading font-bold">🎲 DishDice</h1>
          </Link>

          <nav className="flex items-center space-x-6">
            <Link
              to="/dashboard"
              className="flex items-center space-x-2 hover:text-secondary transition"
            >
              <Home size={20} />
              <span>{t('nav.dashboard')}</span>
            </Link>

            <Link
              to="/shopping-list"
              className="flex items-center space-x-2 hover:text-secondary transition"
            >
              <ShoppingCart size={20} />
              <span>{t('nav.shopping')}</span>
            </Link>

            <Link
              to="/preferences"
              className="flex items-center space-x-2 hover:text-secondary transition"
            >
              <Settings size={20} />
              <span>{t('nav.preferences')}</span>
            </Link>

            <div className="flex items-center space-x-4 border-l border-white/30 pl-6">
              <div className="flex items-center space-x-2">
                <User size={20} />
                <span className="text-sm">{user?.email}</span>
              </div>

              <button
                onClick={handleLogout}
                className="flex items-center space-x-2 bg-white/20 hover:bg-white/30 px-4 py-2 rounded-lg transition"
              >
                <LogOut size={20} />
                <span>{t('nav.logout')}</span>
              </button>
            </div>
          </nav>
        </div>
      </div>
    </header>
  );
};
