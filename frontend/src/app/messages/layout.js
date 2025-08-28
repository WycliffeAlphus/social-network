import { ContactedChats } from "@/components/contacts";

import React from "react";

export default function MessagesLayout({ children }) {
    return (
        <div className="h-screen flex">
            <div className="w-[30%] border-x border-gray-400 overflow-auto scrollbar dark:scrollbar-dark p-2">
                <h2 className="font-semibold mb-7 p-4">Contacts</h2>
                <ContactedChats />
            </div>

            <div className="flex-1 border-r border-gray-400">
                {children}
            </div>
        </div>
    );
}
