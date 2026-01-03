import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ApolloProvider } from '@apollo/client';
import client from './graphql/client';
import AppLayout from './components/layout/AppLayout';
import Dashboard from './pages/Dashboard';
import TeamList from './pages/TeamList';
import DeveloperList from './pages/DeveloperList';

function App() {
  return (
    <ApolloProvider client={client}>
      <BrowserRouter>
        <AppLayout>
          <Routes>
            <Route path="/" element={<Navigate to="/dashboard" replace />} />
            <Route path="/dashboard" element={<Dashboard />} />
            <Route path="/teams" element={<TeamList />} />
            <Route path="/developers" element={<DeveloperList />} />
          </Routes>
        </AppLayout>
      </BrowserRouter>
    </ApolloProvider>
  );
}

export default App;
