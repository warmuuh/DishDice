import React, { useState, useEffect } from 'react';
import { useTranslation } from 'react-i18next';
import { Header } from '../components/Header';
import { userService } from '../services/userService';
import toast from 'react-hot-toast';
import { Save } from 'lucide-react';

export const Preferences: React.FC = () => {
  const { t, i18n } = useTranslation();
  const [preferences, setPreferences] = useState('');
  const [language, setLanguage] = useState('');
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    loadPreferences();
  }, []);

  const loadPreferences = async () => {
    try {
      const data = await userService.getPreferences();
      setPreferences(data.preferences || '');
      const userLanguage = data.language || 'en';
      setLanguage(userLanguage);
      // Only change language if it's different from current
      if (i18n.language !== userLanguage) {
        i18n.changeLanguage(userLanguage);
      }
    } catch (error) {
      toast.error(t('preferences.failed'));
    } finally {
      setLoading(false);
    }
  };

  const handleSave = async () => {
    setSaving(true);
    try {
      await userService.updatePreferences(preferences, language);
      i18n.changeLanguage(language);
      localStorage.setItem('language', language);
      toast.success(t('preferences.saved'));
    } catch (error) {
      toast.error(t('preferences.failed'));
    } finally {
      setSaving(false);
    }
  };

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

      <main className="container mx-auto px-4 py-8">
        <div className="max-w-3xl mx-auto">
          <div className="bg-white rounded-2xl shadow-lg p-8">
            <h1 className="text-3xl font-heading font-bold text-primary mb-2">
              {t('preferences.title')}
            </h1>
            <p className="text-gray-600 mb-6">
              {t('preferences.description')}
            </p>

            <div className="space-y-6">
              {/* Language Selection */}
              <div>
                <label htmlFor="language" className="block text-sm font-medium text-gray-700 mb-2">
                  {t('preferences.language')}
                </label>
                <select
                  id="language"
                  value={language}
                  onChange={(e) => setLanguage(e.target.value)}
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
                  disabled={saving}
                >
                  <option value="en">English</option>
                  <option value="de">Deutsch</option>
                </select>
              </div>

              {/* Food Preferences */}
              <div>
                <label htmlFor="preferences" className="block text-sm font-medium text-gray-700 mb-2">
                  {t('preferences.title')}
                </label>
                <textarea
                  id="preferences"
                  value={preferences}
                  onChange={(e) => setPreferences(e.target.value)}
                  placeholder={t('preferences.placeholder')}
                  className="w-full h-64 px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent resize-none"
                  disabled={saving}
                />
              </div>

              <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                <h3 className="font-semibold text-blue-900 mb-2">{t('preferences.tips.title')}</h3>
                <ul className="text-sm text-blue-800 space-y-1">
                  <li>• {t('preferences.tips.restrictions')}</li>
                  <li>• {t('preferences.tips.cuisines')}</li>
                  <li>• {t('preferences.tips.allergies')}</li>
                  <li>• {t('preferences.tips.cookingStyle')}</li>
                </ul>
              </div>

              <button
                onClick={handleSave}
                disabled={saving}
                className="w-full bg-gradient-to-r from-primary to-accent text-white py-3 px-6 rounded-lg font-semibold hover:shadow-lg transform hover:scale-105 transition disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center space-x-2"
              >
                <Save size={20} />
                <span>{saving ? t('preferences.saving') : t('preferences.save')}</span>
              </button>
            </div>
          </div>
        </div>
      </main>
    </div>
  );
};
