import { ContactedChats } from "@/components/contacts";

import React from "react";

export default function MessagesLayout({ children }) {
    return (
        <div className="bg-red-900 h-screen flex">
            <div className="bg-blue-900 opacity-50 w-[30%] border-x border-gray-200 overflow-auto scrollbar p-4">
                <h2 className="font-semibold mb-12">Contacts</h2>
                <ContactedChats />
            </div>

            <div className="bg-green-900 opacity-50 flex-1 px-4">
                {children}
            </div>
        </div>
    );
}
