"use client";

import {
  createContext,
  useContext,
  useState,
  useCallback,
  useRef,
  useEffect,
} from "react";
import Notification from "@/app/components/common/Notification";

const NotificationContext = createContext();

export const NotificationProvider = ({ children }) => {
  const [notifications, setNotifications] = useState([]);
  const timersRef = useRef({});

  useEffect(() => {
    console.log("Updated notifications:", notifications);
  }, [notifications]);

  const hideNotification = useCallback((id) => {
    setNotifications((prev) =>
      prev.map((n) => (n.id === id ? { ...n, visible: false } : n))
    );
    if (timersRef.current[id]) {
      clearTimeout(timersRef.current[id]);
      delete timersRef.current[id];
    }
  }, []);

  const addNotification = useCallback(
    (notification) => {
      const id = Date.now() + Math.random().toString(36).substr(2, 9);
      const newNotification = {
        id,
        type: notification.type || "info",
        message: notification.message,
        duration: notification.duration ?? 3000,
        position: notification.position || "top-right",
        createdAt: Date.now(),
        link: notification.link || null,
        visible: true, // this controls whether it's showing
      };

      setNotifications((prev) => [newNotification, ...prev]);

      if (newNotification.duration !== 0) {
        const timer = setTimeout(() => {
          hideNotification(id);
        }, newNotification.duration);
        timersRef.current[id] = timer;
      }

      return id;
    },
    [hideNotification]
  );

  const removeNotification = useCallback((id) => {
    setNotifications((prev) => prev.filter((n) => n.id !== id));
    if (timersRef.current[id]) {
      clearTimeout(timersRef.current[id]);
      delete timersRef.current[id];
    }
  }, []);

  const showSuccess = useCallback(
    (msg, opts = {}) =>
      addNotification({ type: "success", message: msg, ...opts }),
    [addNotification]
  );
  const showError = useCallback(
    (msg, opts = {}) =>
      addNotification({ type: "error", message: msg, ...opts }),
    [addNotification]
  );
  const showWarning = useCallback(
    (msg, opts = {}) =>
      addNotification({ type: "warning", message: msg, ...opts }),
    [addNotification]
  );
  const showInfo = useCallback(
    (msg, opts = {}) =>
      addNotification({ type: "info", message: msg, ...opts }),
    [addNotification]
  );

  const clearAllNotifications = () => {
    setNotifications([]);
  };

  return (
    <NotificationContext.Provider
      value={{
        addNotification,
        removeNotification,
        showSuccess,
        showError,
        showWarning,
        showInfo,
        notifications, // full list including invisible
        clearAllNotifications,
      }}
    >
      {children}
      <div className="notifications-container">
        {notifications
          .filter((n) => n.visible) // only show those that are visible
          .map((notification) => (
            <Notification
              key={notification.id}
              type={notification.type}
              message={notification.message}
              duration={notification.duration}
              position={notification.position}
              onClose={() => hideNotification(notification.id)}
            />
          ))}
      </div>
    </NotificationContext.Provider>
  );
};

export function useNotification() {
  const context = useContext(NotificationContext);
  if (!context) {
    throw new Error(
      "useNotification must be used within a NotificationProvider"
    );
  }
  return context;
}
