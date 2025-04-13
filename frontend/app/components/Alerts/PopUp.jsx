// AlertContext.js
import { createContext, useContext, useState } from "react";
import styles from "@/styles/PopUp.module.css";

const AlertContext = createContext(undefined);

export const AlertProvider = ({ children }) => {
  const [alertConfig, setAlertConfig] = useState({
    isOpen: false,
    type: "",
    message: "",
    action: undefined,
  });

  const showAlert = ({ type, message, action }) => {
    setAlertConfig({ isOpen: true, type, message, action });
  };

  const closeAlert = () => {
    setAlertConfig((prev) => ({ ...prev, isOpen: false }));
  };

  return (
    <AlertContext.Provider value={{ alertConfig, showAlert, closeAlert }}>
      {children}
    </AlertContext.Provider>
  );
};

export const useAlert = () => {
  const context = useContext(AlertContext);
  if (!context) {
    throw new Error("useAlert must be used within an AlertProvider");
  }
  return context;
};

// ConfirmAction.js
export const ConfirmAction = () => {
  const { alertConfig, closeAlert } = useAlert();
  if (!alertConfig.isOpen || alertConfig.type !== "confirm") return null;

  const handleConfirm = () => {
    if (alertConfig.action) {
      alertConfig.action();
    }
    closeAlert();
  };

  return (
    <div className={styles.overlay}>
      <div className={styles.popup}>
        <p>{alertConfig.message}</p>
        <div className={styles.buttons}>
          <button className={styles.confirm} onClick={handleConfirm}>
            Continue
          </button>
          <button className={styles.cancel} onClick={closeAlert}>
            Cancel
          </button>
        </div>
      </div>
    </div>
  );
};

export const PopUp = () => {
  const { alertConfig, closeAlert } = useAlert();
  if (!alertConfig.isOpen || alertConfig.type === "confirm") return null;

  const popUpClass =
    alertConfig.type === "success" ? styles.success : styles.failure;

  return (
    <div className={styles.overlay}>
      <div className={`${styles.popup} ${popUpClass}`}>
        <p>{alertConfig.message}</p>
        <button className={styles.close} onClick={closeAlert}>
          Close
        </button>
      </div>
    </div>
  );
};
