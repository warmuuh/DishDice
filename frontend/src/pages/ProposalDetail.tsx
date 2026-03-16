import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { Header } from '../components/Header';
import { DayCard } from '../components/DayCard';
import { RegenerateModal } from '../components/RegenerateModal';
import { proposalService } from '../services/proposalService';
import { mealService } from '../services/mealService';
import type { WeeklyProposal, DailyMealOption } from '../types';
import toast from 'react-hot-toast';
import { Calendar, ArrowLeft } from 'lucide-react';

export const ProposalDetail: React.FC = () => {
  const { t } = useTranslation();
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [proposal, setProposal] = useState<WeeklyProposal | null>(null);
  const [loading, setLoading] = useState(true);
  const [regenerateModalOpen, setRegenerateModalOpen] = useState(false);
  const [regeneratingMealId, setRegeneratingMealId] = useState<string | null>(null);
  const [mealOptions, setMealOptions] = useState<DailyMealOption[]>([]);
  const [optionsLoading, setOptionsLoading] = useState(false);

  useEffect(() => {
    if (id) {
      loadProposal();
    }
  }, [id]);

  const loadProposal = async () => {
    if (!id) return;

    try {
      const data = await proposalService.getProposal(id);
      setProposal(data);
    } catch (error) {
      toast.error(t('proposal.loadFailed'));
      navigate('/dashboard');
    } finally {
      setLoading(false);
    }
  };

  const handleRegenerate = async (mealId: string) => {
    setRegeneratingMealId(mealId);
    setRegenerateModalOpen(true);
    setOptionsLoading(true);

    try {
      const response = await mealService.regenerateMeal(mealId);
      setMealOptions(response.options);
    } catch (error) {
      toast.error(t('meal.regenerateFailed'));
      setRegenerateModalOpen(false);
    } finally {
      setOptionsLoading(false);
    }
  };

  const handleSelectOption = async (optionIndex: number) => {
    if (!regeneratingMealId) return;

    const selectedMeal = mealOptions[optionIndex];
    if (!selectedMeal) return;

    try {
      await mealService.selectOption(regeneratingMealId, {
        option_index: optionIndex,
        menu_name: selectedMeal.menu_name,
        recipe: selectedMeal.recipe,
        shopping_items: selectedMeal.shopping_items,
      });
      toast.success(t('meal.updateSuccess'));
      setRegenerateModalOpen(false);
      setRegeneratingMealId(null);
      loadProposal();
    } catch (error) {
      toast.error(t('meal.updateFailed'));
    }
  };

  const handleAddToShopping = async (mealId: string) => {
    try {
      await mealService.addToShoppingList(mealId);
      toast.success(t('meal.addedToShopping'));
    } catch (error) {
      toast.error(t('meal.addToShoppingFailed'));
    }
  };

  const handleAddAllToShoppingList = async () => {
    if (!id) return;

    try {
      await proposalService.addProposalToShoppingList(id);
      toast.success(t('proposal.allAddedToShopping'));
    } catch (error) {
      toast.error(t('proposal.addToShoppingFailed'));
    }
  };

  const formatWeekRange = (startDate: string) => {
    const start = new Date(startDate);
    const end = new Date(start);
    end.setDate(start.getDate() + 6);

    return `${start.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })} - ${end.toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}`;
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

  if (!proposal) {
    return (
      <div className="min-h-screen bg-gray-50">
        <Header />
        <div className="container mx-auto px-4 py-8 text-center">
          <p className="text-gray-600">{t('proposal.notFound')}</p>
        </div>
      </div>
    );
  }

  const regeneratingMeal = proposal.daily_meals.find((m) => m.id === regeneratingMealId);
  const dayNames = [
    t('days.monday'),
    t('days.tuesday'),
    t('days.wednesday'),
    t('days.thursday'),
    t('days.friday'),
    t('days.saturday'),
    t('days.sunday')
  ];

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />

      <main className="container mx-auto px-4 py-8">
        <button
          onClick={() => navigate('/dashboard')}
          className="flex items-center space-x-2 text-primary hover:text-accent transition mb-6"
        >
          <ArrowLeft size={20} />
          <span>{t('proposal.backToDashboard')}</span>
        </button>

        <div className="bg-white rounded-2xl shadow-lg p-8 mb-8">
          <div className="flex items-center justify-between mb-4">
            <div className="flex items-center space-x-4">
              <Calendar className="text-primary" size={32} />
              <div>
                <h1 className="text-3xl font-heading font-bold text-gray-900">
                  {t('proposal.weeklyPlan')}
                </h1>
                <p className="text-lg text-gray-600">
                  {formatWeekRange(proposal.week_start_date)}
                </p>
              </div>
            </div>
            <button
              onClick={handleAddAllToShoppingList}
              className="bg-gradient-to-r from-success to-success/80 text-white py-3 px-6 rounded-lg font-semibold hover:shadow-lg transform hover:scale-105 transition flex items-center space-x-2"
            >
              <span>{t('proposal.addAllToShopping')}</span>
            </button>
          </div>

          {proposal.week_preferences && (
            <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mt-4">
              <h3 className="font-semibold text-blue-900 mb-1">{t('proposal.weekPreferencesLabel')}</h3>
              <p className="text-blue-800">{proposal.week_preferences}</p>
            </div>
          )}

          {proposal.current_resources && (
            <div className="bg-green-50 border border-green-200 rounded-lg p-4 mt-4">
              <h3 className="font-semibold text-green-900 mb-1">{t('proposal.availableIngredientsLabel')}</h3>
              <p className="text-green-800">{proposal.current_resources}</p>
            </div>
          )}
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
          {proposal.daily_meals
            .sort((a, b) => a.day_of_week - b.day_of_week)
            .map((meal) => (
              <DayCard
                key={meal.id}
                meal={meal}
                dayName={dayNames[meal.day_of_week]}
                onRegenerate={() => handleRegenerate(meal.id)}
                onAddToShopping={() => handleAddToShopping(meal.id)}
              />
            ))}
        </div>
      </main>

      <RegenerateModal
        isOpen={regenerateModalOpen}
        onClose={() => {
          setRegenerateModalOpen(false);
          setRegeneratingMealId(null);
        }}
        options={mealOptions}
        loading={optionsLoading}
        onSelect={handleSelectOption}
        dayName={regeneratingMeal ? dayNames[regeneratingMeal.day_of_week] : ''}
      />
    </div>
  );
};
