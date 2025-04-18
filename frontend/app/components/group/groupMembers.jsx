export default function GroupMembers({handleButtonClick2, isMembersListVisible ,handleButtonClick,isUsersListVisible,members }) {
return (
    <div>
         {/* displaying group memebers */}
    {/* Button to toggle visibility of members list */}
    <div className="flex justify-center items-center gap-4 mt-4">
    <button
        onClick={handleButtonClick2}
        className="mt-4 bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
    >
        {isMembersListVisible ? "Hide Members" : "Show Members"}
    </button>
     {/* Button to toggle the visibility of the users list */}
     <button
            onClick={handleButtonClick}
            className="mt-4 bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
        >
            {isUsersListVisible ? "Hide Users" : "Show People"}
        </button>
    </div>
    {/* Conditionally render the members list if visible */}
    {isMembersListVisible && (
        <div className="mt-4 p-4 bg-blue-300 rounded-lg shadow-md border border-blue-300">
        {members && members.length > 0 ? (
            members.map((member, index) => (
            <div key={index} className="member text-white mb-2">
                <p>{member.nickname}</p>
            </div>
            ))
        ) : (
            <p className="text-gray-400">No members found</p>
        )}
        </div>
    )}
      {/* end of displaying group members */}
    </div>
  

);

}