import React, { useState } from "react";
import { useNotification } from "@/context/NotificationContext";
import { useRouter } from "next/navigation";
import { BiBell } from "react-icons/bi"; 

const NotificationButton = () => {
  const [isOpen, setIsOpen] = useState(false);
  const { notifications, removeNotification, clearAllNotifications } =
    useNotification();
  const router = useRouter();

  const handleNotificationClick = (notification) => {
    if (notification.link) {
      router.push(notification.link);
    }
    removeNotification(notification.id);
    setIsOpen(false);
  };

  return (
    <div className="relative">
      <button
        onClick={() => setIsOpen(!isOpen)}
        className="p-2 rounded-full hover:bg-gray-100 relative"
      >
        <BiBell className="w-6 h-6" />
        {notifications.length > 0 && (
          <span className="absolute -top-1 -right-1 bg-red-500 text-white rounded-full w-5 h-5 text-xs flex items-center justify-center">
            {notifications.length}
          </span>
        )}
      </button>

      {isOpen && (
        <div className="absolute right-0 mt-2 w-80 bg-white rounded-lg shadow-lg border border-gray-200 max-h-96 overflow-y-auto z-50">
          {notifications.length > 0 ? (
            <>
              <div className="flex justify-between items-center px-4 py-2 border-b border-gray-100">
                <span className="text-sm font-medium text-gray-700">
                  Notifications
                </span>
                <button
                  onClick={() => {
                    clearAllNotifications();
                    setIsOpen(false);
                  }}
                  className="text-xs text-red-600 hover:text-red-800"
                >
                  Clear All
                </button>
              </div>
              <div className="py-2">
                {notifications.map((notification) => (
                  <div
                    key={notification.id}
                    onClick={() => handleNotificationClick(notification)}
                    className="px-4 py-3 hover:bg-gray-50 cursor-pointer border-b border-gray-100 last:border-0"
                  >
                    <div className="flex items-start">
                      <div className="flex-1">
                        <p className="text-sm text-gray-800">
                          {notification.message}
                        </p>
                        <p className="text-xs text-gray-500 mt-1">
                          {new Date(
                            notification.createdAt
                          ).toLocaleTimeString()}
                        </p>
                      </div>
                      <div
                        className={`w-2 h-2 rounded-full mt-2 ${getNotificationTypeColor(
                          notification.type
                        )}`}
                      />
                    </div>
                  </div>
                ))}
              </div>
            </>
          ) : (
            <div className="px-4 py-6 text-center text-gray-500">
              No notifications
            </div>
          )}
        </div>
      )}
    </div>
  );
};

// Helper function to get notification color based on type
const getNotificationTypeColor = (type) => {
  switch (type) {
    case "success":
      return "bg-green-500";
    case "error":
      return "bg-red-500";
    case "warning":
      return "bg-yellow-500";
    default:
      return "bg-blue-500";
  }
};

export default NotificationButton;
