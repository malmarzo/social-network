export default function GroupChat() {
     const [selectedUsers, setSelectedUsers] = useState([]);  // Store selected users here
     const createUsersList = async () => {
    if (selectedUsers.length == 0){
        alert("you need at least to invite one person");
        return;
      }
        if (selectedUsers.length > 0) {
                      await sendInvitations(response.group.id, response.group.creator_id, selectedUsers);
                  }
                };
    return (
        <div>
            <UsersList 
                            selectedUsers={selectedUsers} 
                            setSelectedUsers={setSelectedUsers} 
                        />
                         <button onClick={createUsersList}></button>
        </div>
       
    );
}
