import React from 'react'
import Link from 'next/link';

const AuthButton = ({text, href}) => {
    return (
      <button>
        <Link
          href={href}
          className="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-md shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors duration-200"
        >
          {text}
        </Link>
      </button>
    );
}

export default AuthButton
