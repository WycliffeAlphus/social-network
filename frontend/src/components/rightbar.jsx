"use client"

import FollowSuggestion from "./followsuggestions"

export default function Rightbar() {
  return (
    <div className="sticky pt-[3rem] top-0 w-[30%] lg:block hidden max-w-xs h-fit overflow-y-auto">
      <FollowSuggestion />
    </div>
  );
}
