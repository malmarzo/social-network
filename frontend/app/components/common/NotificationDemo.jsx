'use client';

import { useState } from 'react';
import { useNotification } from '@/context/NotificationContext';

/**
 * A demo component to showcase the different types of notifications
 */
export default function NotificationDemo() {
  const { showSuccess, showError, showInfo, showWarning } = useNotification();
  const [message, setMessage] = useState('This is a notification message');
  const [duration, setDuration] = useState(5000);
  const [position, setPosition] = useState('top-right');

  const handleShowNotification = (type) => {
    switch (type) {
      case 'success':
        showSuccess(message, { duration, position });
        break;
      case 'error':
        showError(message, { duration, position });
        break;
      case 'warning':
        showWarning(message, { duration, position });
        break;
      case 'info':
        showInfo(message, { duration, position });
        break;
      default:
        break;
    }
  };

  return (
    <div className="p-6 bg-white rounded-lg shadow-md">
      <h2 className="text-xl font-semibold mb-4">Notification Demo</h2>
      
      <div className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Message
          </label>
          <input
            type="text"
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
        
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Duration (ms)
          </label>
          <input
            type="number"
            value={duration}
            onChange={(e) => setDuration(Number(e.target.value))}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <p className="text-xs text-gray-500 mt-1">Set to 0 for persistent notification</p>
        </div>
        
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Position
          </label>
          <select
            value={position}
            onChange={(e) => setPosition(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="top-right">Top Right</option>
            <option value="top-left">Top Left</option>
            <option value="bottom-right">Bottom Right</option>
            <option value="bottom-left">Bottom Left</option>
            <option value="top-center">Top Center</option>
            <option value="bottom-center">Bottom Center</option>
          </select>
        </div>
        
        <div className="flex flex-wrap gap-2 pt-2">
          <button
            onClick={() => handleShowNotification('success')}
            className="px-4 py-2 bg-green-500 text-white rounded-md hover:bg-green-600 focus:outline-none focus:ring-2 focus:ring-green-500"
          >
            Success
          </button>
          <button
            onClick={() => handleShowNotification('error')}
            className="px-4 py-2 bg-red-500 text-white rounded-md hover:bg-red-600 focus:outline-none focus:ring-2 focus:ring-red-500"
          >
            Error
          </button>
          <button
            onClick={() => handleShowNotification('warning')}
            className="px-4 py-2 bg-yellow-500 text-white rounded-md hover:bg-yellow-600 focus:outline-none focus:ring-2 focus:ring-yellow-500"
          >
            Warning
          </button>
          <button
            onClick={() => handleShowNotification('info')}
            className="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            Info
          </button>
        </div>
      </div>
    </div>
  );
}
