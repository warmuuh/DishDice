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
