import React from "react";

const PostLoader = () => {
  return (
    <div className="w-full bg-white rounded-xl border border-gray-200 p-6 mb-6 shadow-sm">
      <div className="flex justify-between items-start mb-6">
        <div className="animate-pulse bg-gray-200 h-8 w-1/3 rounded-lg" />
        <div className="flex flex-col items-end gap-2">
          <div className="animate-pulse bg-gray-200 h-4 w-24 rounded-lg" />
          <div className="animate-pulse bg-gray-200 h-3 w-20 rounded-lg" />
        </div>
      </div>

      <div className="animate-pulse bg-gray-200 w-full h-[300px] rounded-xl mb-6" />

      <div className="space-y-3 mb-6">
        <div className="animate-pulse bg-gray-200 w-full h-4 rounded-lg" />
        <div className="animate-pulse bg-gray-200 w-4/5 h-4 rounded-lg" />
      </div>

      <div className="flex gap-6 pt-4 border-t border-gray-100">
        <div className="animate-pulse bg-gray-200 w-16 h-6 rounded-lg" />
        <div className="animate-pulse bg-gray-200 w-16 h-6 rounded-lg" />
        <div className="animate-pulse bg-gray-200 w-16 h-6 rounded-lg" />
      </div>
    </div>
  );
};

export default PostLoader;
