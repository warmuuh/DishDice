import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Header } from '../components/Header';
import { proposalService } from '../services/proposalService';
import type { WeeklyProposal } from '../types';
import toast from 'react-hot-toast';
import { Plus, Calendar, Trash2 } from 'lucide-react';

export const Dashboard: React.FC = () => {
  const { t } = useTranslation();
  const [proposals, setProposals] = useState<WeeklyProposal[]>([]);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    loadProposals();
  }, []);

  const loadProposals = async () => {
    try {
      const data = await proposalService.getProposals();
      setProposals(data || []);
    } catch (error) {
      toast.error(t('dashboard.loadFailed'));
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: string, e: React.MouseEvent) => {
    e.stopPropagation();
    if (!confirm(t('dashboard.deleteConfirm'))) {
      return;
    }

    try {
      await proposalService.deleteProposal(id);
      toast.success(t('dashboard.deleteSuccess'));
      setProposals(proposals.filter((p) => p.id !== id));
    } catch (error) {
      toast.error(t('dashboard.deleteFailed'));
    }
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' });
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
        <div className="flex items-center justify-between mb-8">
          <div>
            <h1 className="text-4xl font-heading font-bold text-gray-900 mb-2">
              {t('dashboard.title')}
            </h1>
            <p className="text-gray-600">{t('dashboard.subtitle')}</p>
          </div>

          <button
            onClick={() => navigate('/proposals/new')}
            className="bg-gradient-to-r from-primary to-accent text-white py-3 px-6 rounded-lg font-semibold hover:shadow-lg transform hover:scale-105 transition flex items-center space-x-2"
          >
            <Plus size={20} />
            <span>{t('dashboard.newProposal')}</span>
          </button>
        </div>

        {proposals.length === 0 ? (
          <div className="bg-white rounded-2xl shadow-lg p-12 text-center">
            <div className="text-6xl mb-4">🎲</div>
            <h2 className="text-2xl font-heading font-bold text-gray-900 mb-2">
              {t('dashboard.noPlans')}
            </h2>
            <p className="text-gray-600 mb-6">
              {t('dashboard.getStarted')}
            </p>
            <button
              onClick={() => navigate('/proposals/new')}
              className="bg-gradient-to-r from-primary to-accent text-white py-3 px-6 rounded-lg font-semibold hover:shadow-lg transform hover:scale-105 transition inline-flex items-center space-x-2"
            >
              <Plus size={20} />
              <span>{t('dashboard.createFirst')}</span>
            </button>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {proposals.map((proposal) => (
              <div
                key={proposal.id}
                onClick={() => navigate(`/proposals/${proposal.id}`)}
                className="bg-white rounded-xl shadow-md hover:shadow-xl transition cursor-pointer overflow-hidden group"
              >
                <div className="bg-gradient-to-r from-primary to-accent p-4">
                  <div className="flex items-center justify-between text-white">
                    <div className="flex items-center space-x-2">
                      <Calendar size={20} />
                      <span className="font-semibold">
                        {t('dashboard.weekOf')} {formatDate(proposal.week_start_date)}
                      </span>
                    </div>
                    <button
                      onClick={(e) => handleDelete(proposal.id, e)}
                      className="opacity-0 group-hover:opacity-100 transition hover:bg-white/20 p-2 rounded"
                    >
                      <Trash2 size={18} />
                    </button>
                  </div>
                </div>

                <div className="p-4">
                  <p className="text-sm text-gray-500 mb-2">
                    {t('dashboard.created')} {formatDate(proposal.created_at)}
                  </p>
                  {proposal.week_preferences && (
                    <p className="text-sm text-gray-700 line-clamp-2">
                      {proposal.week_preferences}
                    </p>
                  )}
                </div>
              </div>
            ))}
          </div>
        )}
      </main>
    </div>
  );
};
