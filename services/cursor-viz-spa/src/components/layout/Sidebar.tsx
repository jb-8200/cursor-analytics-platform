import React from 'react';
import { NavLink } from 'react-router-dom';

const Sidebar: React.FC = () => {
  const navItems = [
    { name: 'Dashboard', path: '/dashboard' },
    { name: 'Teams', path: '/teams' },
    { name: 'Developers', path: '/developers' },
  ];

  return (
    <nav
      role="navigation"
      className="hidden md:block w-64 bg-gray-50 border-r border-gray-200 h-screen sticky top-16"
      aria-label="Main navigation"
    >
      <div className="px-4 py-6">
        <ul className="space-y-2">
          {navItems.map((item) => (
            <li key={item.path}>
              <NavLink
                to={item.path}
                className={({ isActive }) =>
                  `block px-4 py-2 rounded-md text-sm font-medium transition-colors ${
                    isActive
                      ? 'bg-primary-50 text-primary-700'
                      : 'text-gray-700 hover:bg-gray-100 hover:text-gray-900'
                  }`
                }
              >
                {item.name}
              </NavLink>
            </li>
          ))}
        </ul>
      </div>
    </nav>
  );
};

export default Sidebar;
