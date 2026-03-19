import React, { useState, useEffect } from 'react';
import { Link, useNavigate, useSearchParams } from 'react-router-dom';
import { useTranslation } from 'react-i18next';
import { useAuth } from '../context/AuthContext';
import { toast } from 'react-toastify';
import { ticketService } from '../services/ticketService';
import { LoadingSpinner } from '../components/LoadingSpinner';

export const Register: React.FC = () => {
  const { t } = useTranslation();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const { register } = useAuth();
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const ticket = searchParams.get('ticket');
  const [validatingTicket, setValidatingTicket] = useState(!!ticket);
  const [ticketValid, setTicketValid] = useState(false);
  const [ticketMessage, setTicketMessage] = useState('');

  useEffect(() => {
    if (ticket) {
      validateTicket(ticket);
    }
  }, [ticket]);

  const validateTicket = async (token: string) => {
    try {
      const response = await ticketService.validateTicket(token);
      if (response.valid) {
        setTicketValid(true);
        toast.success('Valid registration link! Your account will be auto-approved.');
      } else {
        setTicketValid(false);
        setTicketMessage(response.message || 'Invalid ticket');
        toast.error(response.message || 'Invalid registration link');
        setTimeout(() => navigate('/login?error=invalid_ticket'), 3000);
      }
    } catch (error) {
      toast.error('Failed to validate registration link');
      setTimeout(() => navigate('/login?error=ticket_validation_failed'), 3000);
    } finally {
      setValidatingTicket(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!email || !password || !confirmPassword) {
      toast.error(t('validation.required'));
      return;
    }

    if (password !== confirmPassword) {
      toast.error(t('validation.passwordMismatch'));
      return;
    }

    if (password.length < 6) {
      toast.error(t('validation.passwordLength'));
      return;
    }

    setLoading(true);

    try {
      const result = await register({ email, password, ticket: ticket || undefined });

      if (result.status === 'pending') {
        toast.success('Registration successful! Waiting for admin approval.');
        navigate('/waiting-approval');
      } else {
        toast.success(t('auth.registerSuccess'));
        navigate('/dashboard');
      }
    } catch (error: any) {
      toast.error(error.response?.data || t('auth.registerFailed'));
    } finally {
      setLoading(false);
    }
  };

  if (validatingTicket) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-secondary via-accent to-primary flex items-center justify-center p-4">
        <div className="bg-white rounded-2xl shadow-2xl p-8 w-full max-w-md text-center">
          <LoadingSpinner />
          <p className="mt-4 text-gray-600">Validating registration link...</p>
        </div>
      </div>
    );
  }

  if (ticket && !ticketValid) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-secondary via-accent to-primary flex items-center justify-center p-4">
        <div className="bg-white rounded-2xl shadow-2xl p-8 w-full max-w-md text-center">
          <div className="text-red-500 text-5xl mb-4">❌</div>
          <h2 className="text-2xl font-heading font-bold text-gray-900 mb-2">
            Invalid Registration Link
          </h2>
          <p className="text-gray-600 mb-6">{ticketMessage}</p>
          <p className="text-sm text-gray-500">Redirecting to login...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-secondary via-accent to-primary flex items-center justify-center p-4">
      <div className="bg-white rounded-2xl shadow-2xl p-8 w-full max-w-md">
        {ticketValid && (
          <div className="mb-4 p-3 bg-green-100 border border-green-400 rounded-lg">
            <p className="text-green-800 text-sm font-semibold">
              ✅ Using pre-approved registration link
            </p>
          </div>
        )}

        <div className="text-center mb-8">
          <h1 className="text-4xl font-heading font-bold text-primary mb-2">🎲 {t('app.name')}</h1>
          <p className="text-gray-600">{t('auth.createAccount')}</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label htmlFor="email" className="block text-sm font-medium text-gray-700 mb-2">
              {t('auth.email')}
            </label>
            <input
              id="email"
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
              placeholder={t('auth.emailPlaceholder')}
              disabled={loading}
            />
          </div>

          <div>
            <label htmlFor="password" className="block text-sm font-medium text-gray-700 mb-2">
              {t('auth.password')}
            </label>
            <input
              id="password"
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
              placeholder="••••••••"
              disabled={loading}
            />
          </div>

          <div>
            <label htmlFor="confirmPassword" className="block text-sm font-medium text-gray-700 mb-2">
              {t('auth.confirmPassword')}
            </label>
            <input
              id="confirmPassword"
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              className="w-full px-4 py-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
              placeholder="••••••••"
              disabled={loading}
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full bg-gradient-to-r from-secondary to-accent text-white py-3 px-4 rounded-lg font-semibold hover:shadow-lg transform hover:scale-105 transition disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? t('auth.creatingAccount') : t('auth.signUp')}
          </button>
        </form>

        <div className="mt-6 text-center">
          <p className="text-gray-600">
            {t('auth.hasAccount')}{' '}
            <Link to="/login" className="text-primary font-semibold hover:text-accent transition">
              {t('auth.login')}
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
};
