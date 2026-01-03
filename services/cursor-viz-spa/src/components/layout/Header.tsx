import React from 'react';

const Header: React.FC = () => {
  return (
    <header
      role="banner"
      className="sticky top-0 z-50 bg-white border-b border-gray-200 shadow-sm"
    >
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex items-center justify-between h-16">
          <div className="flex items-center">
            <h1 className="text-2xl font-bold text-gray-900">
              Cursor Analytics Dashboard
            </h1>
          </div>

          <div className="flex items-center space-x-4">
            {/* User menu placeholder */}
            <div className="user-menu">
              <button
                type="button"
                className="flex items-center text-sm rounded-full focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary-500"
                aria-label="User menu"
              >
                <span className="sr-only">Open user menu</span>
                <div className="h-8 w-8 rounded-full bg-gray-300 flex items-center justify-center">
                  <span className="text-gray-600 font-medium">U</span>
                </div>
              </button>
            </div>
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;
