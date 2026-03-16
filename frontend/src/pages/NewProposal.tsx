import React, { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Header } from '../components/Header';
import { proposalService } from '../services/proposalService';
import toast from 'react-hot-toast';
import { Wand2, Info } from 'lucide-react';

export const NewProposal: React.FC = () => {
  const { t } = useTranslation();
  const [weekStartDate, setWeekStartDate] = useState('');
  const [weekPreferences, setWeekPreferences] = useState('');
  const [currentResources, setCurrentResources] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  // Get next Monday as default
  const getNextMonday = () => {
    const today = new Date();
    const dayOfWeek = today.getDay();
    const daysUntilMonday = dayOfWeek === 0 ? 1 : (8 - dayOfWeek) % 7;
    const nextMonday = new Date(today);
    nextMonday.setDate(today.getDate() + daysUntilMonday);
    return nextMonday.toISOString().split('T')[0];
  };

  React.useEffect(() => {
    setWeekStartDate(getNextMonday());
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!weekStartDate) {
      toast.error(t('proposal.selectDate'));
      return;
    }

    setLoading(true);

    try {
      const proposal = await proposalService.createProposal({
        week_start_date: weekStartDate,
        week_preferences: weekPreferences || undefined,
        current_resources: currentResources || undefined,
      });

      toast.success(t('proposal.generated'));
      navigate(`/proposals/${proposal.id}`);
    } catch (error: any) {
      toast.error(error.response?.data || t('proposal.generateFailed'));
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />

      <main className="container mx-auto px-4 py-8">
        <div className="max-w-3xl mx-auto">
          <div className="bg-white rounded-2xl shadow-lg p-8">
            <h1 className="text-3xl font-heading font-bold text-primary mb-2">
              {t('proposal.createNew')}
            </h1>
            <p className="text-gray-600 mb-6">
              {t('proposal.aiGenerate')}
            </p>

            <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-6 flex items-start gap-3">
              <Info className="text-blue-600 flex-shrink-0 mt-0.5" size={20} />
              <p className="text-sm text-blue-800">
                {t('proposal.preferencesHint')}{' '}
                <Link to="/preferences" className="font-semibold underline hover:text-blue-900">
                  {t('preferences.title')}
                </Link>
              </p>
            </div>

            <form onSubmit={handleSubmit} className="space-y-6">
              <div>
                <label htmlFor="weekStartDate" className="block text-sm font-medium text-gray-700 mb-2">
                  {t('proposal.weekStarting')}
                </label>
                <input
                  id="weekStartDate"
                  type="date"
                  value={weekStartDate}
                  onChange={(e) => setWeekStartDate(e.target.value)}
                  className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
                  disabled={loading}
                  required
                />
              </div>

              <div>
                <label htmlFor="weekPreferences" className="block text-sm font-medium text-gray-700 mb-2">
                  {t('proposal.weekPreferences')}
                </label>
                <textarea
                  id="weekPreferences"
                  value={weekPreferences}
                  onChange={(e) => setWeekPreferences(e.target.value)}
                  placeholder={t('proposal.weekPreferencesPlaceholder')}
                  className="w-full h-32 px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent resize-none"
                  disabled={loading}
                />
              </div>

              <div>
                <label htmlFor="currentResources" className="block text-sm font-medium text-gray-700 mb-2">
                  {t('proposal.availableIngredients')}
                </label>
                <textarea
                  id="currentResources"
                  value={currentResources}
                  onChange={(e) => setCurrentResources(e.target.value)}
                  placeholder={t('proposal.availableIngredientsPlaceholder')}
                  className="w-full h-32 px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent resize-none"
                  disabled={loading}
                />
                <p className="text-sm text-gray-500 mt-1">
                  {t('proposal.ingredientsHint')}
                </p>
              </div>

              {loading && (
                <div className="bg-primary/10 border border-primary/20 rounded-lg p-6 text-center">
                  <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mx-auto mb-4"></div>
                  <p className="text-primary font-semibold">
                    {t('proposal.generating')}
                  </p>
                  <p className="text-sm text-gray-600 mt-2">
                    {t('proposal.generatingTime')}
                  </p>
                </div>
              )}

              <div className="flex space-x-4">
                <button
                  type="button"
                  onClick={() => navigate('/dashboard')}
                  className="flex-1 bg-gray-200 text-gray-700 py-3 px-6 rounded-lg font-semibold hover:bg-gray-300 transition"
                  disabled={loading}
                >
                  {t('proposal.cancel')}
                </button>

                <button
                  type="submit"
                  disabled={loading}
                  className="flex-1 bg-gradient-to-r from-primary to-accent text-white py-3 px-6 rounded-lg font-semibold hover:shadow-lg transform hover:scale-105 transition disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center space-x-2"
                >
                  <Wand2 size={20} />
                  <span>{loading ? t('proposal.generatingBtn') : t('proposal.generate')}</span>
                </button>
              </div>
            </form>
          </div>
        </div>
      </main>
    </div>
  );
};
