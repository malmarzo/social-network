'use client';

import NotificationDemo from '../components/common/NotificationDemo';

export default function NotificationDemoPage() {
  return (
    <div className="container mx-auto py-8 px-4">
      <h1 className="text-2xl font-bold mb-6">Notification System Demo</h1>
      <p className="mb-6 text-gray-700">
        This page demonstrates the reusable notification system that can be used throughout the application,
        including in the chat functionality.
      </p>
      
      <NotificationDemo />
      
      <div className="mt-8 p-4 bg-gray-100 rounded-lg">
        <h2 className="text-xl font-semibold mb-2">About the Notification System</h2>
        <p className="mb-4">
          The notification system provides a consistent way to display feedback to users across the application.
          It supports different types of notifications (success, error, warning, info) and various positioning options.
        </p>
        <h3 className="text-lg font-medium mb-2">Features:</h3>
        <ul className="list-disc pl-5 space-y-1">
          <li>Multiple notification types with appropriate styling</li>
          <li>Configurable duration (including persistent notifications)</li>
          <li>Flexible positioning around the screen</li>
          <li>Auto-dismiss functionality</li>
          <li>Manual dismiss option</li>
          <li>Stacking of multiple notifications</li>
        </ul>
        <p className="mt-4">
          This system is integrated with the chat functionality to provide feedback for message sending,
          connection status, and new message notifications.
        </p>
      </div>
    </div>
  );
}
