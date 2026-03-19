import React, { useEffect, useState } from 'react';
import { Header } from '../components/Header';
import { LoadingSpinner } from '../components/LoadingSpinner';
import { adminService } from '../services/adminService';
import type { AdminUser } from '../types';
import { toast } from 'react-toastify';
import { Check, X } from 'lucide-react';

export const AdminPanel: React.FC = () => {
  const [users, setUsers] = useState<AdminUser[]>([]);
  const [filter, setFilter] = useState<'all' | 'pending' | 'approved' | 'rejected'>('pending');
  const [loading, setLoading] = useState(true);
  const [showTicketModal, setShowTicketModal] = useState(false);
  const [generatedLink, setGeneratedLink] = useState('');
  const [creatingTicket, setCreatingTicket] = useState(false);

  useEffect(() => {
    loadUsers();
  }, [filter]);

  const loadUsers = async () => {
    setLoading(true);
    try {
      if (filter === 'pending') {
        const data = await adminService.getPendingUsers();
        setUsers(data || []);
      } else if (filter === 'all') {
        const data = await adminService.getAllUsers();
        setUsers(data || []);
      } else {
        const allUsers = await adminService.getAllUsers();
        setUsers(allUsers ? allUsers.filter(u => u.status === filter) : []);
      }
    } catch (error) {
      toast.error('Failed to load users');
      setUsers([]);
    } finally {
      setLoading(false);
    }
  };

  const handleApprove = async (userId: string) => {
    try {
      await adminService.approveUser(userId);
      toast.success('User approved successfully!');
      loadUsers();
    } catch (error) {
      toast.error('Failed to approve user');
    }
  };

  const handleReject = async (userId: string) => {
    if (!confirm('Are you sure you want to reject this user?')) return;

    try {
      await adminService.rejectUser(userId);
      toast.success('User rejected');
      loadUsers();
    } catch (error) {
      toast.error('Failed to reject user');
    }
  };

  const handleCreateTicket = async () => {
    setCreatingTicket(true);
    try {
      const response = await adminService.createTicket();
      setGeneratedLink(response.registration_link);
      setShowTicketModal(true);
      toast.success('Registration link created!');
    } catch (error) {
      toast.error('Failed to create registration link');
    } finally {
      setCreatingTicket(false);
    }
  };

  const copyToClipboard = () => {
    navigator.clipboard.writeText(generatedLink);
    toast.success('Link copied to clipboard!');
  };

  const getStatusBadge = (status: string) => {
    const colors = {
      pending: 'bg-yellow-100 text-yellow-800',
      approved: 'bg-green-100 text-green-800',
      rejected: 'bg-red-100 text-red-800',
    };
    return (
      <span className={`px-3 py-1 rounded-full text-sm font-semibold ${colors[status as keyof typeof colors]}`}>
        {status}
      </span>
    );
  };

  return (
    <div className="min-h-screen bg-gray-50">
      <Header />

      <main className="container mx-auto px-4 py-8">
        <h1 className="text-4xl font-heading font-bold text-gray-900 mb-8">
          Admin Panel
        </h1>

        {/* Create Registration Link Button */}
        <div className="mb-6">
          <button
            onClick={handleCreateTicket}
            disabled={creatingTicket}
            className="px-6 py-3 bg-gradient-to-r from-primary to-accent text-white rounded-lg font-semibold hover:shadow-lg transform hover:scale-105 transition disabled:opacity-50"
          >
            {creatingTicket ? 'Generating...' : 'Create Registration Link'}
          </button>
        </div>

        {/* Ticket Modal */}
        {showTicketModal && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
            <div className="bg-white rounded-xl shadow-2xl p-6 max-w-lg w-full">
              <h3 className="text-2xl font-heading font-bold text-gray-900 mb-4">
                Registration Link Created
              </h3>
              <p className="text-gray-600 mb-4">
                Share this link with the new user. It's valid for 2 weeks and can only be used once.
              </p>
              <div className="bg-gray-100 p-4 rounded-lg mb-4 break-all">
                <code className="text-sm">{generatedLink}</code>
              </div>
              <div className="flex gap-3">
                <button
                  onClick={copyToClipboard}
                  className="flex-1 px-4 py-2 bg-primary text-white rounded-lg hover:bg-opacity-90 transition"
                >
                  Copy Link
                </button>
                <button
                  onClick={() => setShowTicketModal(false)}
                  className="px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition"
                >
                  Close
                </button>
              </div>
            </div>
          </div>
        )}

        {/* Filter Tabs */}
        <div className="mb-6 flex space-x-2">
          {['all', 'pending', 'approved', 'rejected'].map((f) => (
            <button
              key={f}
              onClick={() => setFilter(f as any)}
              className={`px-4 py-2 rounded-lg font-semibold transition ${
                filter === f
                  ? 'bg-gradient-to-r from-primary to-accent text-white'
                  : 'bg-white text-gray-700 hover:bg-gray-100'
              }`}
            >
              {f.charAt(0).toUpperCase() + f.slice(1)}
            </button>
          ))}
        </div>

        {/* User List */}
        <div className="bg-white rounded-xl shadow-md overflow-hidden">
          <div className="bg-gradient-to-r from-primary to-accent p-6">
            <h2 className="text-2xl font-heading font-bold text-white">
              User Management
            </h2>
          </div>

          {loading ? (
            <div className="p-12 text-center">
              <LoadingSpinner />
            </div>
          ) : !users || users.length === 0 ? (
            <div className="p-12 text-center text-gray-600">
              No {filter !== 'all' ? filter : ''} users found
            </div>
          ) : (
            <div className="divide-y">
              {users.map((user) => (
                <div key={user.id} className="p-4 flex items-center justify-between hover:bg-gray-50">
                  <div className="flex-1">
                    <p className="font-semibold text-gray-900">{user.email}</p>
                    <p className="text-sm text-gray-500">
                      Registered: {new Date(user.created_at).toLocaleString()}
                    </p>
                    <p className="text-sm text-gray-500">
                      Role: <span className="font-semibold">{user.role}</span>
                    </p>
                  </div>

                  <div className="flex items-center gap-4">
                    {getStatusBadge(user.status)}

                    {user.status === 'pending' && (
                      <div className="flex gap-2">
                        <button
                          onClick={() => handleApprove(user.id)}
                          className="flex items-center gap-1 px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition"
                        >
                          <Check size={16} />
                          Approve
                        </button>
                        <button
                          onClick={() => handleReject(user.id)}
                          className="flex items-center gap-1 px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition"
                        >
                          <X size={16} />
                          Reject
                        </button>
                      </div>
                    )}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </main>
    </div>
  );
};
