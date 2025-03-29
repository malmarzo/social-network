// import { useState } from 'react';

// const EventNotificationCard = ({ content, onDismiss }) => {
//   const [isVisible, setIsVisible] = useState(true);

//   // Handle click on the "OK" button to hide the card
//   const handleDismiss = () => {
//     setIsVisible(false);
//     if (onDismiss) {
//       onDismiss(); // Call onDismiss callback if provided
//     }
//   };

//   return (
//     isVisible && (
//       <div className="card">
//         <div className="card-content">
//           <h3>New Event Notification</h3>
//           <p>{`"${content}" `}</p>
//         </div>
//         <button className="btn" onClick={handleDismiss}>
//           OK
//         </button>
//       </div>
//     )
//   );
// };

// export default EventNotificationCard;

import { useState } from 'react';

const EventNotificationCard = ({ content, onDismiss }) => {
  const [isVisible, setIsVisible] = useState(true);

  // Handle click on the "OK" button to hide the card
  const handleDismiss = () => {
    setIsVisible(false);
    if (onDismiss) {
      onDismiss(); // Call onDismiss callback if provided
    }
  };

  // Inline styles for the card
  const cardStyle = {
    backgroundColor: '#ffffff',
    borderRadius: '8px',
    boxShadow: '0 4px 8px rgba(0, 0, 0, 0.1)',
    width: '300px',
    padding: '20px',
    position: 'fixed',
    top: '20px',
    right: '20px',
    zIndex: 1000,
    fontFamily: 'Arial, sans-serif',
  };

  // Inline styles for the button
  const buttonStyle = {
    backgroundColor: '#007bff',
    color: 'white',
    border: 'none',
    padding: '10px 20px',
    borderRadius: '25px',
    fontSize: '1rem',
    cursor: 'pointer',
    transition: 'background-color 0.3s ease',
    marginTop: '15px',
  };

  // Button hover and active styles
  const buttonHoverStyle = {
    backgroundColor: '#0056b3',
  };

  const buttonActiveStyle = {
    backgroundColor: '#003f7d',
    transform: 'translateY(2px)',
  };

  // Card content style
  const cardContentStyle = {
    textAlign: 'center',
  };

  const handleMouseEnter = (e) => {
    e.target.style.backgroundColor = buttonHoverStyle.backgroundColor;
  };

  const handleMouseLeave = (e) => {
    e.target.style.backgroundColor = buttonStyle.backgroundColor;
  };

  const handleMouseDown = (e) => {
    e.target.style.backgroundColor = buttonActiveStyle.backgroundColor;
    e.target.style.transform = 'translateY(2px)';
  };

  const handleMouseUp = (e) => {
    e.target.style.backgroundColor = buttonHoverStyle.backgroundColor;
    e.target.style.transform = 'none';
  };

  return (
    isVisible && (
      <div style={cardStyle}>
        <div style={cardContentStyle}>
          <h3>New Event Notification</h3>
          <p>{`"${content}"`}</p>
        </div>
        <button
          style={buttonStyle}
          onClick={handleDismiss}
          onMouseEnter={handleMouseEnter}
          onMouseLeave={handleMouseLeave}
          onMouseDown={handleMouseDown}
          onMouseUp={handleMouseUp}
        >
          OK
        </button>
      </div>
    )
  );
};

export default EventNotificationCard;

