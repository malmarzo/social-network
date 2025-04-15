import { useState } from 'react';

const EventNotificationCard = ({ content, onDismiss }) => {
  const [isVisible, setIsVisible] = useState(true);

  
  const handleDismiss = () => {
    setIsVisible(false);
    if (onDismiss) {
      onDismiss(); 
    }
  };


  return (
    isVisible && (
      <div style={styles.cardStyle}>
        <div style={styles.cardContentStyle}>
          <h3>New Event Notification</h3>
          <p>{`"${content}"`}</p>
        </div>
        <button
          style={styles.buttonStyle}
          onClick={handleDismiss}
         
        >
          Discard
        </button>
      </div>
    )
  );
};

export default EventNotificationCard;

const styles = {
  cardStyle: {
    backgroundColor: '#fdfdfd',
    borderRadius: '16px',
    boxShadow: '0 6px 20px rgba(0, 0, 0, 0.08)',
    width: '100%',
    padding: '24px',
    marginBottom: '16px',
    fontFamily: `'Segoe UI', Tahoma, Geneva, Verdana, sans-serif`,
    boxSizing: 'border-box',
    transition: 'transform 0.2s ease, box-shadow 0.3s ease',
  },

  cardHoverStyle: {
    transform: 'scale(1.02)',
    boxShadow: '0 10px 24px rgba(0, 0, 0, 0.12)',
  },

  cardContentStyle: {
    textAlign: 'left',
    color: '#333',
    fontSize: '1.05rem',
    lineHeight: '1.6',
  },

  titleStyle: {
    fontSize: '1.5rem',
    fontWeight: '600',
    marginBottom: '12px',
    color: '#111',
  },

  subtitleStyle: {
    fontSize: '1rem',
    color: '#555',
    marginBottom: '16px',
  },

  buttonStyle: {
    backgroundColor: '#4f46e5', // elegant indigo
    color: '#ffffff',
    border: 'none',
    padding: '10px 24px',
    borderRadius: '30px',
    fontSize: '1rem',
    cursor: 'pointer',
    transition: 'background-color 0.3s ease, transform 0.2s ease',
  },

  buttonHoverStyle: {
    backgroundColor: '#4338ca',
    transform: 'translateY(-2px)',
  },
};
