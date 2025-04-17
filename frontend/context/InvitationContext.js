// src/contexts/InvitationContext.js
import { createContext, useState } from "react";

export const InvitationContext = createContext();

export const InvitationProvider = ({ children }) => {
  const [invitations, setInvitations] = useState([]);

  return (
    <InvitationContext.Provider value={{ invitations, setInvitations }}>
      {children}
    </InvitationContext.Provider>
  );
};
