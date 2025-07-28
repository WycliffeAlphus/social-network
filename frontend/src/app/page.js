"use client"

import FollowSuggestion from "../components/followsuggestions";
import Rightbar from "../components/rightbar";

export default function Home() {
  return (
    <div className="flex min-h-screen">
      <main className="flex-1 border-x mr-[20px] border-gray-400">
        <div className="lg:hidden">
          <FollowSuggestion />
        </div>
        <div className="p-4 border-t lg:border-0 border-gray-400">
          <div>Lorem Ipsum dummy content</div>
          <div>Lorem Ipsum dummy content</div>
          <div>Lorem Ipsum dummy content</div>
          <div>Lorem Ipsum dummy content</div>
          <div>Lorem Ipsum dummy content</div>
          <div>Lorem Ipsum dummy content</div>
          <div>Lorem Ipsum dummy content</div>
          <div>Lorem Ipsum dummy content</div>
          <div>Lorem Ipsum dummy content</div>
          <div>Lorem Ipsum dummy content</div>
        </div>
      </main>
      <Rightbar />
    </div>
  );
}
