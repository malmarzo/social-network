// sendInvitations.js
import { invokeAPI } from "@/utils/invokeAPI";

export const sendInvitations = async (groupId, invitedBy, users) => {
    await Promise.all(
        users.map((userId) =>
            invokeAPI("groups/invite", {
                group_id: groupId,
                user_id: userId,
                invited_by: invitedBy,
            }, "POST")
        )
    );
};
