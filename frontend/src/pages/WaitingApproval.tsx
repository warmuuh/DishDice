import React from 'react';
import { Link } from 'react-router-dom';

export const WaitingApproval: React.FC = () => {
  return (
    <div className="min-h-screen bg-gradient-to-br from-primary via-accent to-secondary flex items-center justify-center p-4">
      <div className="bg-white rounded-2xl shadow-2xl p-8 w-full max-w-md text-center">
        <div className="text-6xl mb-4">⏳</div>
        <h1 className="text-3xl font-heading font-bold text-primary mb-4">
          Waiting for Approval
        </h1>
        <p className="text-gray-600 mb-6">
          Your registration was successful! An administrator will review your account soon.
          You'll be able to log in once approved.
        </p>
        <Link
          to="/login"
          className="inline-block bg-gradient-to-r from-primary to-accent text-white py-3 px-6 rounded-lg font-semibold hover:shadow-lg transform hover:scale-105 transition"
        >
          Back to Login
        </Link>
      </div>
    </div>
  );
};
