import React from "react";
import ProfileCard from "./ProfileCard";

const ProfileWindow = () => {
  return (
    <div>
          <h1
          style={{fontSize:"1.8rem", fontWeight:"bold"}}
          >Profile</h1>
      <ProfileCard />
    </div>
  );
};

export default ProfileWindow;
