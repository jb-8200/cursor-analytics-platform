import React from 'react';
import Header from './Header';
import Sidebar from './Sidebar';

interface AppLayoutProps {
  children: React.ReactNode;
}

const AppLayout: React.FC<AppLayoutProps> = ({ children }) => {
  return (
    <div className="flex min-h-screen flex-col">
      <Header />
      <div className="flex flex-1">
        <Sidebar />
        <main className="flex-1 px-4 py-6 sm:px-6 lg:px-8 bg-gray-100">
          {children}
        </main>
      </div>
    </div>
  );
};

export default AppLayout;
